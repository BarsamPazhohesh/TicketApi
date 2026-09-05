package otp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"ticket-api/internal/config"
	"ticket-api/internal/db/sms_type_messages_relation"
	"ticket-api/internal/db/sms_warehouse"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/services/cache"
	"ticket-api/internal/services/otp"
	"ticket-api/internal/services/smsloader"
	"ticket-api/internal/testutil"
	"time"
)

func TestOTPService_Flow(t *testing.T) {
	db := testutil.SetupTestSQLDB(t)
	defer db.Close()

	// Mock SMS Gateway server
	smsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.FormValue("userName") == "" || r.FormValue("passWord") == "" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("-2\n<html>error</html>"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("3257258978\n<html>success</html>"))
	}))
	defer smsServer.Close()

	_ = os.Setenv("OTP_URL", smsServer.URL)
	_ = os.Setenv("OTP_USERNAME", "test_user")
	_ = os.Setenv("OTP_PASSWORD", "test_pass")
	_ = os.Setenv("OTP_SENDER_NUMBER", "1000")

	root := testutil.FindProjectRoot()
	_ = config.Load(filepath.Join(root, "config.yaml"))

	cacheSvc := cache.NewCacheService(nil) // in-memory/nil redis safe
	warehouseRepo := repository.NewSMSWarehouseRepository(sms_warehouse.New(db))
	smsLoader := smsloader.NewSMSLoaderService(sms_type_messages_relation.New(db))
	otpSvc := otp.NewOTPService(warehouseRepo, cacheSvc, smsLoader)

	ctx := context.Background()
	phone := "09121234567"

	t.Run("SendOTP and verify warehouse record creation", func(t *testing.T) {
		res, apiErr := otpSvc.SendOTP(ctx, phone)
		if apiErr != nil {
			t.Fatalf("SendOTP returned error: %v", apiErr)
		}
		if res.Message != "کد تایید با موفقیت ارسال شد" {
			t.Fatalf("expected Persian message 'کد تایید با موفقیت ارسال شد', got: %s", res.Message)
		}

		// Wait briefly for async goroutine to execute
		time.Sleep(100 * time.Millisecond)

		// Check warehouse entry
		pending, apiErr := warehouseRepo.GetPending(ctx)
		if apiErr != nil {
			t.Fatalf("GetPending failed: %v", apiErr)
		}

		// Since mock SMS server responded with 200, status should have been updated to Sent (1)
		record, apiErr := warehouseRepo.GetByID(ctx, 1)
		if apiErr != nil {
			t.Fatalf("GetByID failed: %v", apiErr)
		}

		if record.ReceiverPhoneNumber != phone {
			t.Fatalf("expected phone %s, got %s", phone, record.ReceiverPhoneNumber)
		}

		if record.Status != dto.SMSStatusSent {
			t.Fatalf("expected status Sent (%d), got %d", dto.SMSStatusSent, record.Status)
		}

		_ = pending
	})

	t.Run("GenerateRandomOTP produces 6-digit numeric code", func(t *testing.T) {
		code, err := otpSvc.GenerateRandomOTP()
		if err != nil {
			t.Fatalf("GenerateRandomOTP failed: %v", err)
		}
		if len(code) != 6 {
			t.Fatalf("expected 6 digits, got length %d: %s", len(code), code)
		}
	})

	t.Run("VerifyOTP without cache returns ErrOTPExpired", func(t *testing.T) {
		_, apiErr := otpSvc.VerifyOTP(ctx, "09129999999", "123456")
		if apiErr == nil {
			t.Fatal("expected error, got nil")
		}
		if apiErr.Err.Code != errx.ErrOTPExpired {
			t.Fatalf("expected ErrOTPExpired, got code %d", apiErr.Err.Code)
		}
	})
}
