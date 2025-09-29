package captcha

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"strconv"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/env"
	"ticket-api/internal/errx"
	"time"

	"github.com/dchest/captcha"
	"github.com/patrickmn/go-cache"
)

// CaptchaService wraps a go-cache for captchas
type CaptchaService struct {
	cache *cache.Cache
}

// NewCaptchaService creates a captcha service with TTL
func NewCaptchaService() *CaptchaService {
	cfg := config.Get().Captcha

	return &CaptchaService{
		cache: cache.New(
			time.Duration(cfg.TimeoutMinutes)*time.Minute,
			time.Duration(cfg.CleanupInterval)*time.Minute,
		),
	}
}

// GenerateCaptcha creates a new captcha and stores the answer in cache
func (s *CaptchaService) GenerateCaptcha() (*dto.CaptchaResultDTO, error) {
	captchaID := captcha.NewLen(6)
	digits := captcha.RandomDigits(6)

	img := captcha.NewImage(captchaID, digits, 240, 80)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	s.cache.Set(captchaID, digits, time.Duration(config.Get().Captcha.TimeoutMinutes)*time.Minute)

	result := &dto.CaptchaResultDTO{
		ID:    captchaID,
		Image: imgBase64,
	}

	// if gin is in debug mode, include the answer for easier testing
	if env.GetEnvString("GIN_MODE", "debug") == "debug" {
		answer := ""
		for _, d := range digits {
			answer += strconv.Itoa(int(d))
		}
		result.Answer = answer
	}

	return result, nil
}

// VerifyCaptcha checks if the given captcha ID and answer are correct
func (s *CaptchaService) VerifyCaptcha(id, answer string) *errx.APIError {
	val, found := s.cache.Get(id)
	if !found {
		return errx.Respond(errx.ErrExpiredCaptcha, errors.New("captcha expired or not found"))
	}

	expected := val.([]byte)

	s.cache.Delete(id)

	if !compareCaptcha(expected, answer) {
		return errx.Respond(errx.ErrIncorrectCaptcha, errors.New("incorrect captcha"))
	}

	return nil // success
}

func compareCaptcha(b []byte, s string) bool {
	if len(b) != len(s) {
		return false
	}

	for i := range b {
		if b[i] != s[i]-'0' {
			return false
		}
	}
	return true
}
