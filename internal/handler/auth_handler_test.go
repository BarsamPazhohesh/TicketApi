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

func TestAuthHandler_Endpoints(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	// 1. SignUp with Password (requires admin/agent auth cookie)
	t.Run("POST /api/v1/auth/SignUp/ without auth cookie -> 401", func(t *testing.T) {
		body, _ := json.Marshal(dto.SignUpWithPasswordDTO{
			Username:     "09121113344",
			Password:     "password123",
			DepartmentID: 1,
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/SignUp/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 Unauthorized, got %d", w.Code)
		}
	})

	t.Run("POST /api/v1/auth/SignUp/ with valid admin auth cookie -> 200/201", func(t *testing.T) {
		adminToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 1, "admin", []int64{1})

		body, _ := json.Marshal(dto.SignUpWithPasswordDTO{
			Username:     "09121113344",
			Password:     "password123",
			DepartmentID: 1,
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/SignUp/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Auth.CookieName,
			Value: adminToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK && w.Code != http.StatusCreated {
			t.Fatalf("expected 200/201 on valid signup, got status %d, body: %s", w.Code, w.Body.String())
		}
	})

	// 2. Login with Password -> Sets Auth Cookie
	t.Run("POST /api/v1/auth/Login/ with valid password -> sets cookie", func(t *testing.T) {
		body, _ := json.Marshal(dto.LoginWithPasswordDTO{
			Username: "09121113344",
			Password: "password123",
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/Login/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on valid login, got %d, body: %s", w.Code, w.Body.String())
		}

		cookies := w.Result().Cookies()
		var foundAuth bool
		for _, ck := range cookies {
			if ck.Name == config.Get().Auth.CookieName {
				foundAuth = true
				break
			}
		}
		if !foundAuth {
			t.Fatal("expected auth cookie in response")
		}
	})

	// 3. Login with No Auth (Admin/Agent provisioning)
	t.Run("POST /api/v1/auth/LoginWithNoAuth/ without auth cookie -> 401", func(t *testing.T) {
		body, _ := json.Marshal(dto.LoginWitNoAuthDTO{
			Username:     "09125556677",
			DepartmentID: 1,
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/LoginWithNoAuth/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 on unauthenticated LoginWithNoAuth, got %d", w.Code)
		}
	})

	t.Run("POST /api/v1/auth/LoginWithNoAuth/ with admin auth cookie -> 200/201", func(t *testing.T) {
		adminToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 1, "admin", []int64{1})

		body, _ := json.Marshal(dto.LoginWitNoAuthDTO{
			Username:     "09125556677",
			DepartmentID: 1,
		})

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/LoginWithNoAuth/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Auth.CookieName,
			Value: adminToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusCreated && w.Code != http.StatusOK {
			t.Fatalf("expected 200/201 on NoAuth login, got %d", w.Code)
		}
	})

	// 4. GET /api/v1/auth/CheckToken/ tests
	t.Run("GET /api/v1/auth/CheckToken/ without cookies -> valid: false, tokenType: none", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/CheckToken/", nil)
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on CheckToken, got %d", w.Code)
		}

		var resp dto.CheckTokenResponseDTO
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode CheckToken response: %v", err)
		}

		if resp.Valid || resp.TokenType != "none" || resp.PhoneVerified {
			t.Fatalf("expected valid: false, tokenType: none, phoneVerified: false, got: %+v", resp)
		}
	})

	t.Run("GET /api/v1/auth/CheckToken/ with valid captcha cookie (no phone) -> tokenType: captcha", func(t *testing.T) {
		captchaToken := testutil.GenerateTestCaptchaTokenWithPhone(t, app.Services.Token, "", "")

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/CheckToken/", nil)
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Captcha.CookieName,
			Value: captchaToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on CheckToken, got %d", w.Code)
		}

		var resp dto.CheckTokenResponseDTO
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode CheckToken response: %v", err)
		}

		if !resp.Valid || resp.TokenType != "captcha" || resp.PhoneVerified {
			t.Fatalf("expected valid: true, tokenType: captcha, phoneVerified: false, got: %+v", resp)
		}
	})

	t.Run("GET /api/v1/auth/CheckToken/ with verified phone captcha cookie -> tokenType: guest", func(t *testing.T) {
		guestPhone := "09129998877"
		captchaToken := testutil.GenerateTestCaptchaTokenWithPhone(t, app.Services.Token, "", guestPhone)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/CheckToken/", nil)
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Captcha.CookieName,
			Value: captchaToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on CheckToken, got %d", w.Code)
		}

		var resp dto.CheckTokenResponseDTO
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode CheckToken response: %v", err)
		}

		if !resp.Valid || resp.TokenType != "guest" || !resp.PhoneVerified || resp.PhoneNumber != guestPhone {
			t.Fatalf("expected valid: true, tokenType: guest, phoneVerified: true, phone: %s, got: %+v", guestPhone, resp)
		}
	})

	t.Run("GET /api/v1/auth/CheckToken/ with valid auth cookie -> tokenType: auth", func(t *testing.T) {
		authToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 42, "09121234567", []int64{1, 2})

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/CheckToken/", nil)
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Auth.CookieName,
			Value: authToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on CheckToken, got %d", w.Code)
		}

		var resp dto.CheckTokenResponseDTO
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode CheckToken response: %v", err)
		}

		if !resp.Valid || resp.TokenType != "auth" || !resp.PhoneVerified || resp.UserID == nil || *resp.UserID != 42 || resp.Username != "09121234567" {
			t.Fatalf("expected valid: true, tokenType: auth, userId: 42, username: 09121234567, got: %+v", resp)
		}
	})
}
