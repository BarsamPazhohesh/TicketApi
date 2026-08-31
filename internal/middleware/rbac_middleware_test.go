package middleware_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"ticket-api/internal/db/api_keys"
	"ticket-api/internal/db/api_routes"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/errx"
	"ticket-api/internal/middleware"
	"ticket-api/internal/repository"
	"ticket-api/internal/security"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	errx.NewRegistry(nil)
}

func setupTestRouter(t *testing.T) (*gin.Engine, *security.SecurityRegistry) {
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	migration, err := os.ReadFile("../../cmd/migrate/migrations/000014_create_permissions_table.up.sql")
	if err != nil {
		migration, err = os.ReadFile("../../../cmd/migrate/migrations/000014_create_permissions_table.up.sql")
		if err != nil {
			t.Fatalf("failed to read migration file: %v", err)
		}
	}

	initSQL := `
	CREATE TABLE IF NOT EXISTS roles (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL UNIQUE, status INT2 DEFAULT 1, deleted INT2 DEFAULT 0, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')), deleted_at TEXT DEFAULT NULL);
	CREATE TABLE IF NOT EXISTS api_routes (id INTEGER PRIMARY KEY AUTOINCREMENT, route TEXT NOT NULL, method TEXT NOT NULL, description TEXT, status INT2 DEFAULT 1, deleted INT2 DEFAULT 0, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')), deleted_at TEXT DEFAULT NULL, UNIQUE(route, method));
	CREATE TABLE IF NOT EXISTS api_keys (id INTEGER PRIMARY KEY AUTOINCREMENT, key TEXT NOT NULL UNIQUE, description TEXT, status INT2 DEFAULT 1, deleted INT2 DEFAULT 0);
	`
	db.Exec(initSQL)
	db.Exec(string(migration))

	rolesRelQueries := roles_relations.New(db)
	apiKeysQueries := api_keys.New(db)
	apiRoutesQueries := api_routes.New(db)

	repo := repository.NewRolesRelationRepository(rolesRelQueries, apiKeysQueries, apiRoutesQueries)
	registry := security.NewSecurityRegistry(repo)
	if err := registry.Reload(context.Background()); err != nil {
		t.Fatalf("failed to reload security registry: %v", err)
	}

	router := gin.New()
	return router, registry
}

func TestDynamicRBACGuardMiddleware(t *testing.T) {
	router, registry := setupTestRouter(t)

	// Protected group with RBAC Guard
	v1 := router.Group("/api/v1")
	v1.Use(func(c *gin.Context) {
		// Mock user injection based on header
		role := c.GetHeader("X-Test-Role")
		if role == "admin" {
			c.Set("user", &token.AuthClaims{UserID: 1, RoleIDs: []int64{1}})
		} else if role == "customer" {
			c.Set("user", &token.AuthClaims{UserID: 2, RoleIDs: []int64{3}})
		}
		c.Next()
	})
	v1.Use(middleware.DynamicRBACGuardMiddleware(registry))
	v1.POST("tickets/GetTicketsList/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 1. Unauthenticated / Missing user context -> 401
	req := httptest.NewRequest("POST", "/api/v1/tickets/GetTicketsList/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for missing user context, got %d", w.Code)
	}

	// 2. Customer accessing admin/agent list -> 403 Forbidden
	req = httptest.NewRequest("POST", "/api/v1/tickets/GetTicketsList/", nil)
	req.Header.Set("X-Test-Role", "customer")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for customer, got %d", w.Code)
	}

	// 3. Admin accessing list -> 200 OK
	req = httptest.NewRequest("POST", "/api/v1/tickets/GetTicketsList/", nil)
	req.Header.Set("X-Test-Role", "admin")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK for admin, got %d", w.Code)
	}
}
