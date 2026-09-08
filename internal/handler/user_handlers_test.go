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

func TestUserHandler_Endpoints(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	// Seed user
	u1, err := app.Repos.Users.CreateUserWithPassword(t.Context(), dto.SignUpWithPasswordDTO{
		Username:     "09129998877",
		Password:     "userpass123",
		DepartmentID: 1,
	})
	if err != nil {
		t.Fatalf("failed seeding user: %v", err)
	}

	// 1. Without Auth Token -> 401
	t.Run("POST /api/v1/users/GetUserByID/ without Auth -> 401", func(t *testing.T) {
		body, _ := json.Marshal(dto.IDRequest[int64]{ID: u1.ID})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/GetUserByID/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 Unauthorized, got %d", w.Code)
		}
	})

	// 2. With Admin Auth Token -> 200
	t.Run("POST /api/v1/users/GetUserByID/ with Admin Token -> 200", func(t *testing.T) {
		adminToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 1, "admin", []int64{1})

		body, _ := json.Marshal(dto.IDRequest[int64]{ID: u1.ID})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/GetUserByID/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Auth.CookieName,
			Value: adminToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d, body: %s", w.Code, w.Body.String())
		}
	})

	// 3. GetUserByUsername
	t.Run("POST /api/v1/users/GetUserByUsername/ with Admin Token -> 200", func(t *testing.T) {
		adminToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 1, "admin", []int64{1})

		body, _ := json.Marshal(dto.UsernameDTO{Username: "09129998877"})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/GetUserByUsername/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Auth.CookieName,
			Value: adminToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", w.Code)
		}
	})

	// 4. GetUsersByIDs
	t.Run("POST /api/v1/users/GetUsersByIDs/ with Admin Token -> 200", func(t *testing.T) {
		adminToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 1, "admin", []int64{1})

		body, _ := json.Marshal(dto.UserIDsDTO{IDs: []int64{u1.ID}})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/GetUsersByIDs/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  config.Get().Auth.CookieName,
			Value: adminToken,
		})
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", w.Code)
		}
	})
}
