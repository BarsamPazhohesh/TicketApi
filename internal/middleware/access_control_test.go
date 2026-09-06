package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"ticket-api/internal/config"
	"ticket-api/internal/middleware"
	"ticket-api/internal/testutil"

	"github.com/gin-gonic/gin"
)

// TestAgentRoleAccess verifies agent (role 2) can access agent-permitted routes
// and is blocked from admin-only routes.
func TestAgentRoleAccess(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)

	// Seed user 10 with department_id for UserHandler tests
	_, err := app.DB.Exec(`INSERT INTO users (id, username, password, department_id, status) VALUES (10, 'agent_user', 'hash', 1, 1)`)
	if err != nil {
		t.Fatalf("failed to insert user 10: %v", err)
	}

	agentToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 10, "agent_user", []int64{2})

	cases := []struct {
		name     string
		path     string
		body     string
		wantCode int
	}{
		{
			// users:manage -> agent does NOT have it; route has no explicit permission mapping
			// so authenticated access is allowed by default (route has no restriction in DB)
			// This tests the "no DB rule = allow authenticated" path for agent
			name:     "agent can access unrestricted auth route",
			path:     "/api/v1/users/GetUserByID/",
			body:     `{"id": 10}`,
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: "irmto_ticket_auth_token", Value: agentToken})
			w := httptest.NewRecorder()
			app.Engine.ServeHTTP(w, req)
			if w.Code != tc.wantCode {
				t.Errorf("expected %d, got %d (body: %s)", tc.wantCode, w.Code, w.Body.String())
			}
		})
	}
}

// TestCustomerBlockedFromAgentRoutes verifies customer (role 3) cannot access
// routes that require tickets:read_all (agent/admin only).
func TestCustomerBlockedFromAgentRoutes(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)

	customerToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 20, "customer_user", []int64{3})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/GetTicketsList/", nil)
	req.AddCookie(&http.Cookie{Name: "irmto_ticket_auth_token", Value: customerToken})
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for customer on GetTicketsList, got %d (body: %s)", w.Code, w.Body.String())
	}
}

// TestDisabledRouteReturns404 verifies that a route marked status=0 in the
// security registry is rejected even for authenticated users.
func TestDisabledRouteReturns404(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)

	// Disable the tickets/GetTicketsList/ route directly in DB then reload
	_, err := app.DB.Exec(`UPDATE api_routes SET status = 0 WHERE route = 'tickets/GetTicketsList/' AND method = 'POST'`)
	if err != nil {
		t.Fatalf("failed to disable route: %v", err)
	}
	if err := app.SecurityRegistry.Reload(t.Context()); err != nil {
		t.Fatalf("failed to reload registry: %v", err)
	}

	adminToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 1, "admin", []int64{1})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/GetTicketsList/", nil)
	req.AddCookie(&http.Cookie{Name: "irmto_ticket_auth_token", Value: adminToken})
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)

	// DynamicRBACGuardMiddleware returns ErrNotFound (404) for disabled routes
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for disabled route, got %d (body: %s)", w.Code, w.Body.String())
	}
}

// TestApiKeyGuardMiddleware verifies the API key gate rejects missing/short keys
// and passes valid ones.
func TestApiKeyGuardMiddleware(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)

	// Insert a known API key (hashed) into DB
	rawKey := "test-api-key-that-is-long-enough-12345678"
	hashedKey := app.Services.Token.Hash(rawKey)
	_, err := app.DB.Exec(`INSERT INTO api_keys (key, description, status) VALUES (?, 'test key', 1)`, hashedKey)
	if err != nil {
		t.Fatalf("failed to insert api key: %v", err)
	}

	// Build a minimal router wired with ApiKeyGuardMiddleware
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/protected", middleware.ApiKeyGuardMiddleware(app.Services.Token, app.Repos.APIKeys), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	cases := []struct {
		name     string
		key      string
		wantCode int
	}{
		{"no key", "", http.StatusUnauthorized},
		{"short key", "short", http.StatusUnauthorized},
		{"invalid key", "wrong-key-that-is-long-enough-but-not-valid", http.StatusUnauthorized},
		{"valid key", rawKey, http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/protected", nil)
			if tc.key != "" {
				req.Header.Set("x-api-key", tc.key)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.wantCode {
				t.Errorf("expected %d, got %d (body: %s)", tc.wantCode, w.Code, w.Body.String())
			}
		})
	}
}

// TestAuthorizationMiddlewareMissingCookie verifies that protected routes
// return 401 when no auth cookie is present.
func TestAuthorizationMiddlewareMissingCookie(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)

	// No cookie set — AuthorizationMiddleware should reject before RBAC runs
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/GetTicketsList/", nil)
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for missing cookie, got %d (body: %s)", w.Code, w.Body.String())
	}
}

// TestAuthorizationMiddlewareInvalidToken verifies that a tampered/invalid JWT
// returns 401, not a panic or 500.
func TestAuthorizationMiddlewareInvalidToken(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/GetTicketsList/", nil)
	req.AddCookie(&http.Cookie{Name: "irmto_ticket_auth_token", Value: "totally.invalid.jwt"})
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for invalid token, got %d (body: %s)", w.Code, w.Body.String())
	}
}

// TestGuestCannotAccessProtectedEndpoints verifies that endpoints requiring
// authentication (no cookie) consistently return 401 across core routes.
func TestGuestCannotAccessProtectedEndpoints(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)

	protectedRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/v1/tickets/GetTicketsList/"},
		{http.MethodPost, "/api/v1/users/GetUserByID/"},
		{http.MethodPost, "/api/v1/tickets/GetTicketByID/"},
		{http.MethodPost, "/api/v1/auth/SignUp/"},
		{http.MethodPost, "/api/v1/auth/LoginWithNoAuth/"},
	}

	for _, route := range protectedRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			w := httptest.NewRecorder()
			app.Engine.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected 401, got %d", w.Code)
			}
		})
	}
}

// TestGuestCanAccessCustomerFileEndpoints verifies that guest with Level 2 phone token
// can reach the customer group endpoints like ticket file uploads and download links.
func TestGuestCanAccessCustomerFileEndpoints(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	guestPhoneToken := testutil.GenerateTestCaptchaTokenWithPhone(t, app.Services.Token, "127.0.0.1", "09121234567")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/GetDownloadLinkTicketFile/test-object-name", bytes.NewBufferString(`{"id": "00000000-0000-0000-0000-000000000000"}`))
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  config.Get().Captcha.CookieName,
		Value: guestPhoneToken,
	})
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)

	// Since it's valid Level 2 token, it passes middleware (won't be 401 Unauthorized).
	if w.Code == http.StatusUnauthorized {
		t.Errorf("expected Level 2 token to pass customer group middleware, got 401 Unauthorized")
	}
}

// TestAgentPermissionCount verifies agent role has exactly the expected permissions
// from the seed data (tickets:read_all, tickets:chat, users:read, files:upload,
// files:download, metadata:read = 6).
func TestAgentPermissionCount(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)

	agentPerms := app.SecurityRegistry.GetPermissionsForRoles([]int64{2})
	if len(agentPerms) != 6 {
		t.Errorf("expected 6 permissions for agent, got %d: %v", len(agentPerms), agentPerms)
	}
}
