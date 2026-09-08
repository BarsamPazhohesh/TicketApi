package otp

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/env"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/services/cache"
	"ticket-api/internal/services/smsloader"
	"time"
)

const (
	defaultOTPTTLMinutes         = 2
	defaultRetryIntervalMinutes  = 5
	defaultMaxVerifyAttempts     = 5
	otpCachePrefix               = "otp:"
	otpAttemptsPrefix            = "otp:attempts:"
)

type SMSConfig struct {
	URL          string
	Username     string
	Password     string
	SenderNumber string
}

type OTPService struct {
	warehouseRepo *repository.SMSWarehouseRepository
	cache         *cache.CacheService
	smsLoader     *smsloader.SMSLoaderService
	smsConfig     SMSConfig
	httpClient    *http.Client
}

func NewOTPService(
	warehouseRepo *repository.SMSWarehouseRepository,
	cache *cache.CacheService,
	smsLoader *smsloader.SMSLoaderService,
) *OTPService {
	cfg := SMSConfig{
		URL:          env.GetEnvString("OTP_URL", "http://localhost:8080/sms/send"),
		Username:     env.GetEnvString("OTP_USERNAME", ""),
		Password:     env.GetEnvString("OTP_PASSWORD", ""),
		SenderNumber: env.GetEnvString("OTP_SENDER_NUMBER", ""),
	}

	svc := &OTPService{
		warehouseRepo: warehouseRepo,
		cache:         cache,
		smsLoader:     smsLoader,
		smsConfig:     cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	retryInterval := time.Duration(config.Get().OTP.RetryIntervalMinutes) * time.Minute
	if retryInterval <= 0 {
		retryInterval = defaultRetryIntervalMinutes * time.Minute
	}

	// Start background retry worker for pending SMS records
	go svc.startPendingSMSRetryWorker(retryInterval)

	return svc
}

// GenerateRandomOTP generates a secure 6-digit numeric OTP code
func (s *OTPService) GenerateRandomOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// SendOTP creates OTP, stores in Redis with configurable TTL, logs to warehouse, and asynchronously sends SMS
func (s *OTPService) SendOTP(ctx context.Context, phone string) (*dto.SendOTPResponseDTO, *errx.APIError) {
	code, err := s.GenerateRandomOTP()
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	ttlMinutes := config.Get().OTP.CodeTTLMinutes
	if ttlMinutes <= 0 {
		ttlMinutes = defaultOTPTTLMinutes
	}
	ttl := time.Duration(ttlMinutes) * time.Minute

	// Store code in Redis with TTL
	cacheKey := otpCachePrefix + phone
	if err := s.cache.Set(ctx, cacheKey, code, ttl); err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	// Reset any previous failed attempts counter for this phone
	_ = s.cache.Delete(ctx, otpAttemptsPrefix+phone)

	// Fetch message template from SMSLoader (DB-backed, in-memory cached)
	msgTemplate := "کد تایید شما: %s"
	var smsTypeID int64 = 1
	if s.smsLoader != nil {
		if tpl, ok := s.smsLoader.GetMessageByType("otp"); ok && tpl != "" {
			msgTemplate = tpl
		}
		if tid, ok := s.smsLoader.GetTypeIDByTitle("otp"); ok && tid > 0 {
			smsTypeID = tid
		}
	}
	message := fmt.Sprintf(msgTemplate, code)

	// Log initial record in sms_warehouse with status Pending (0)
	record, apiErr := s.warehouseRepo.Create(ctx, smsTypeID, phone, message, dto.SMSStatusPending)
	if apiErr != nil {
		return nil, apiErr
	}

	// Asynchronous send via goroutine
	go func(recordID int64, targetPhone, msg string) {
		asyncCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		sendErr := s.sendSMSViaHTTP(asyncCtx, targetPhone, msg)
		if sendErr != nil {
			log.Printf("[OTP] Async SMS send failed for phone %s (record ID %d): %v", targetPhone, recordID, sendErr)
			_ = s.warehouseRepo.UpdateStatus(context.Background(), recordID, dto.SMSStatusFailed)
			return
		}

		_ = s.warehouseRepo.UpdateStatus(context.Background(), recordID, dto.SMSStatusSent)
	}(record.ID, phone, message)

	return &dto.SendOTPResponseDTO{
		Message: "کد تایید با موفقیت ارسال شد",
	}, nil
}

// VerifyOTP validates the OTP code against cache
func (s *OTPService) VerifyOTP(ctx context.Context, phone, code string) (*dto.VerifyOTPResponseDTO, *errx.APIError) {
	cacheKey := otpCachePrefix + phone
	var cachedCode string

	found, err := s.cache.Get(ctx, cacheKey, &cachedCode)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	if !found {
		return nil, errx.Respond(errx.ErrOTPExpired, nil)
	}

	if subtle.ConstantTimeCompare([]byte(cachedCode), []byte(code)) != 1 {
		attemptsKey := otpAttemptsPrefix + phone

		ttlMinutes := config.Get().OTP.CodeTTLMinutes
		if ttlMinutes <= 0 {
			ttlMinutes = defaultOTPTTLMinutes
		}
		attempts, _ := s.cache.Incr(ctx, attemptsKey, time.Duration(ttlMinutes)*time.Minute)

		maxAttempts := int64(config.Get().OTP.MaxVerifyAttempts)
		if maxAttempts <= 0 {
			maxAttempts = defaultMaxVerifyAttempts
		}

		if attempts >= maxAttempts {
			// Lockout threshold reached: purge OTP code and attempts counter
			_ = s.cache.Delete(ctx, cacheKey)
			_ = s.cache.Delete(ctx, attemptsKey)
			return nil, errx.Respond(errx.ErrOTPMaxAttemptsExceeded, nil)
		}

		return nil, errx.Respond(errx.ErrOTPInvalid, nil)
	}

	// Consume OTP upon successful verification (single-use)
	_ = s.cache.Delete(ctx, cacheKey)
	_ = s.cache.Delete(ctx, otpAttemptsPrefix+phone)

	return &dto.VerifyOTPResponseDTO{
		Valid:   true,
		Message: "تایید با موفقیت انجام شد",
	}, nil
}

// sendSMSViaHTTP sends the SMS payload to the configured OTP endpoint via form-urlencoded POST
// with required keys: userName, passWord, smsText, reciverNumber, senderNumber
func (s *OTPService) sendSMSViaHTTP(ctx context.Context, reciverNumber, smsText string) error {
	form := url.Values{
		"userName":      {s.smsConfig.Username},
		"passWord":      {s.smsConfig.Password},
		"senderNumber":  {s.smsConfig.SenderNumber},
		"reciverNumber": {reciverNumber},
		"smsText":       {smsText},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.smsConfig.URL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sms gateway http error: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fields := strings.Fields(string(bodyBytes))
	if len(fields) == 0 {
		return fmt.Errorf("empty response from sms gateway")
	}

	resultCode, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid response format from sms gateway: %s", fields[0])
	}

	if resultCode <= 0 {
		return fmt.Errorf("sms gateway rejected send, code: %d", resultCode)
	}

	return nil
}

// startPendingSMSRetryWorker processes pending sms_warehouse entries periodically
func (s *OTPService) startPendingSMSRetryWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		s.retryPendingSMS()
	}
}

// retryPendingSMS retrieves and retries all pending records
func (s *OTPService) retryPendingSMS() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	pendingList, err := s.warehouseRepo.GetPending(ctx)
	if err != nil {
		log.Printf("[OTP Worker] Failed to fetch pending SMS records: %v", err)
		return
	}

	if len(pendingList) == 0 {
		return
	}

	log.Printf("[OTP Worker] Retrying %d pending SMS records", len(pendingList))

	for _, item := range pendingList {
		sendErr := s.sendSMSViaHTTP(ctx, item.ReceiverPhoneNumber, item.Message)
		if sendErr != nil {
			log.Printf("[OTP Worker] Retry failed for record ID %d: %v", item.ID, sendErr)
			_ = s.warehouseRepo.UpdateStatus(ctx, item.ID, dto.SMSStatusFailed)
			continue
		}

		_ = s.warehouseRepo.UpdateStatus(ctx, item.ID, dto.SMSStatusSent)
		log.Printf("[OTP Worker] Successfully resent SMS for record ID %d", item.ID)
	}
}
