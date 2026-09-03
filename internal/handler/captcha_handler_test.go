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

func TestCaptchaHandler_Endpoints(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	t.Run("GET /api/v1/captcha/GetCaptcha/ -> 200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/captcha/GetCaptcha/", nil)
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on captcha get, got %d", w.Code)
		}
	})

	t.Run("POST /api/v1/captcha/VerifyCaptcha/ -> 200 sets cookie", func(t *testing.T) {
		res, err := app.Services.Captcha.GenerateCaptcha()
		if err != nil {
			t.Fatalf("failed generating captcha: %v", err)
		}

		body, _ := json.Marshal(dto.CaptchaVerifyRequest{
			ID:          res.ID,
			Captcha:     res.Answer,
			PhoneNumber: "09123456789",
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/captcha/VerifyCaptcha/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on captcha verify, got %d, body: %s", w.Code, w.Body.String())
		}

		cookies := w.Result().Cookies()
		var foundCaptcha bool
		for _, ck := range cookies {
			if ck.Name == config.Get().Captcha.CookieName {
				foundCaptcha = true
				break
			}
		}
		if !foundCaptcha {
			t.Error("expected captcha cookie to be set")
		}
	})
}
