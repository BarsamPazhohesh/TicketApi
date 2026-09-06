package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/testutil"
)

func TestOTPHandler_Endpoints(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	captchaToken := testutil.GenerateTestCaptchaTokenWithPhone(t, app.Services.Token, "", "")

	t.Run("POST /api/v1/otp/send/ with invalid phone number -> 400", func(t *testing.T) {
		body, _ := json.Marshal(dto.SendOTPDTO{
			PhoneNumber: "invalid-phone",
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/send/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Captcha.CookieName,
			Value: captchaToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d, body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("POST /api/v1/otp/send/ with valid phone number -> 200", func(t *testing.T) {
		body, _ := json.Marshal(dto.SendOTPDTO{
			PhoneNumber: "09121234567",
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/send/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Captcha.CookieName,
			Value: captchaToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d, body: %s", w.Code, w.Body.String())
		}

		var resp dto.SendOTPResponseDTO
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if resp.Message == "" {
			t.Fatal("expected non-empty message in response")
		}
	})

	t.Run("POST /api/v1/otp/verify/ with invalid format -> 400", func(t *testing.T) {
		body, _ := json.Marshal(dto.VerifyOTPDTO{
			PhoneNumber: "09121234567",
			Code:        "123", // too short (needs 6 digits)
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/verify/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Captcha.CookieName,
			Value: captchaToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request for short code, got %d", w.Code)
		}
	})

	t.Run("POST /api/v1/otp/verify/ success sets captcha cookie with phone", func(t *testing.T) {
		phone := "09121234567"
		_, apiErr := app.Services.OTP.SendOTP(t.Context(), phone)
		if apiErr != nil {
			t.Fatalf("failed sending OTP: %v", apiErr)
		}

		var cachedCode string
		found, err := app.Services.Cache.Get(t.Context(), "otp:"+phone, &cachedCode)
		if err != nil || !found {
			t.Fatalf("failed to retrieve cached OTP code: %v", err)
		}

		body, _ := json.Marshal(dto.VerifyOTPDTO{
			PhoneNumber: phone,
			Code:        cachedCode,
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/verify/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Captcha.CookieName,
			Value: captchaToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK for valid verify, got %d, body: %s", w.Code, w.Body.String())
		}

		cookies := w.Result().Cookies()
		var foundCaptchaToken string
		expectedCookieName := config.Get().Captcha.CookieName
		for _, ck := range cookies {
			if ck.Name == expectedCookieName {
				foundCaptchaToken = ck.Value
				break
			}
		}
		if foundCaptchaToken == "" {
			t.Fatalf("expected cookie %s to be set on verify success, got cookies: %+v", expectedCookieName, cookies)
		}

		claims, parseErr := app.Services.Token.ParseCaptchaToken(foundCaptchaToken)
		if parseErr != nil {
			t.Fatalf("failed parsing captcha token: %v", parseErr)
		}
		if claims.PhoneNumber != phone {
			t.Fatalf("expected phone %s in captcha claims, got %s", phone, claims.PhoneNumber)
		}
	})
}
