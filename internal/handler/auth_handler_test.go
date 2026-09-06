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
}
