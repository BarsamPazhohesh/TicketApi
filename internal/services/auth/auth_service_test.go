package auth_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"ticket-api/internal/config"
	"ticket-api/internal/db/api_keys"
	"ticket-api/internal/db/api_routes"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/db/users"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/services/auth"
	"ticket-api/internal/services/token"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	_ = os.Setenv("JWT_SECRET", "this-is-a-very-secure-32-byte-long-secret-key-12345")
}

func setupTestAuthDB(t *testing.T) (*sql.DB, *auth.AuthService) {
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory: %v", err)
	}

	schema := `
	CREATE TABLE departments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		status INT2 NOT NULL DEFAULT 1
	);
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT,
		department_id INTEGER NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		status INT2 NOT NULL DEFAULT 1,
		FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
	);
	CREATE TABLE roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		status INT2 NOT NULL DEFAULT 1
	);
	CREATE TABLE permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		status INT2 NOT NULL DEFAULT 1
	);
	CREATE TABLE users_roles_relation (
		user_id INTEGER NOT NULL,
		role_id INTEGER NOT NULL,
		status INT2 NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		PRIMARY KEY (user_id, role_id)
	);
	CREATE TABLE api_routes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		route TEXT NOT NULL,
		method TEXT NOT NULL,
		description TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		status INT2 NOT NULL DEFAULT 1
	);
	CREATE TABLE api_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		status INT2 NOT NULL DEFAULT 1
	);
	CREATE TABLE api_keys_roles_relation (
		api_key_id INTEGER NOT NULL,
		role_id INTEGER NOT NULL,
		status INT2 NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		PRIMARY KEY (api_key_id, role_id)
	);
	CREATE TABLE ticket_types_roles_relation (
		ticket_type_id INTEGER NOT NULL,
		role_id INTEGER NOT NULL,
		status INT2 NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		PRIMARY KEY (ticket_type_id, role_id)
	);
	CREATE TABLE api_routes_roles_relation (
		api_route_id INTEGER NOT NULL,
		role_id INTEGER NOT NULL,
		status INT2 NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		PRIMARY KEY (api_route_id, role_id)
	);
	CREATE TABLE roles_permissions_relation (
		role_id INTEGER NOT NULL,
		permission_id INTEGER NOT NULL,
		status INT2 NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		PRIMARY KEY (role_id, permission_id)
	);
	CREATE TABLE api_routes_permissions_relation (
		api_route_id INTEGER NOT NULL,
		permission_id INTEGER NOT NULL,
		status INT2 NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		deleted_at TEXT DEFAULT NULL,
		PRIMARY KEY (api_route_id, permission_id)
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to init db schema: %v", err)
	}

	_, _ = db.Exec(`INSERT INTO departments (id, title) VALUES (1, 'Support');`)
	_, _ = db.Exec(`INSERT INTO roles (id, name) VALUES (1, 'admin'), (3, 'customer');`)

	errx.NewRegistry(db)
	config.Load("../../../config.yaml")
	userRepo := repository.NewUsersRepository(users.New(db))
	rolesRelRepo := repository.NewRolesRelationRepository(
		roles_relations.New(db),
		api_keys.New(db),
		api_routes.New(db),
	)
	tokenSvc := token.NewTokenService(nil)

	authSvc := auth.NewAuthService(userRepo, rolesRelRepo, tokenSvc)
	return db, authSvc
}

func TestAuthService_SignUpAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, authSvc := setupTestAuthDB(t)
	defer db.Close()

	ctx := context.Background()

	// 1. Sign Up
	signUpReq := dto.SignUpWithPasswordDTO{
		Username:     "09123456789",
		Password:     "securepassword123",
		DepartmentID: 1,
	}

	res, apiErr := authSvc.SignUpWithPassword(ctx, signUpReq)
	if apiErr != nil {
		t.Fatalf("unexpected signup error: %v", apiErr)
	}
	if res.ID <= 0 {
		t.Fatalf("expected positive user ID, got %d", res.ID)
	}

	// 2. Duplicate SignUp should fail
	_, dupErr := authSvc.SignUpWithPassword(ctx, signUpReq)
	if dupErr == nil {
		t.Fatal("expected conflict error on duplicate signup, got nil")
	}

	// 3. Login with wrong password
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/Login/", nil)

	wrongLogin := dto.LoginWithPasswordDTO{
		Username: "09123456789",
		Password: "wrongpassword",
	}
	loginErr := authSvc.LoginWithPassword(c, wrongLogin)
	if loginErr == nil {
		t.Fatal("expected invalid credentials error on wrong password, got nil")
	}

	// 4. Login with correct password
	correctLogin := dto.LoginWithPasswordDTO{
		Username: "09123456789",
		Password: "securepassword123",
	}
	validLoginErr := authSvc.LoginWithPassword(c, correctLogin)
	if validLoginErr != nil {
		t.Fatalf("unexpected login error: %v", validLoginErr)
	}

	// Check cookie
	cookies := w.Result().Cookies()
	var authCookie *http.Cookie
	for _, ck := range cookies {
		if ck.Name == config.Get().Auth.CookieName {
			authCookie = ck
			break
		}
	}
	if authCookie == nil {
		t.Error("expected auth cookie to be set")
	}
}

func TestAuthService_LoginWithNoAuth(t *testing.T) {
	db, authSvc := setupTestAuthDB(t)
	defer db.Close()

	ctx := context.Background()

	req := dto.LoginWitNoAuthDTO{
		Username:     "09999999999",
		DepartmentID: 1,
	}

	// First call -> creates user (status 201)
	res1, status1, err1 := authSvc.LoginWithNoAuth(ctx, req)
	if err1 != nil || status1 != http.StatusCreated {
		t.Fatalf("expected 201 Created on first no-auth login, got status %d, err %v", status1, err1)
	}

	// Second call -> returns existing user (status 200)
	res2, status2, err2 := authSvc.LoginWithNoAuth(ctx, req)
	if err2 != nil || status2 != http.StatusOK {
		t.Fatalf("expected 200 OK on existing no-auth login, got status %d, err %v", status2, err2)
	}
	if res1.ID != res2.ID {
		t.Errorf("expected same user ID across no-auth logins, got %d and %d", res1.ID, res2.ID)
	}
}
