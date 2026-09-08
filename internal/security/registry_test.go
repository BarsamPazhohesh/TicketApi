package security_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"ticket-api/internal/db/api_keys"
	"ticket-api/internal/db/api_routes"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/repository"
	"ticket-api/internal/security"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, *security.SecurityRegistry) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	migration, err := os.ReadFile("../../cmd/migrate/migrations/000014_create_permissions_table.up.sql")
	if err != nil {
		// Try root-relative path if running from package dir
		migration, err = os.ReadFile("../../../cmd/migrate/migrations/000014_create_permissions_table.up.sql")
		if err != nil {
			t.Fatalf("failed to read migration file: %v", err)
		}
	}

	// Schema dependencies
	initSQL := `
	CREATE TABLE IF NOT EXISTS roles (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL UNIQUE, status INT2 DEFAULT 1, deleted INT2 DEFAULT 0, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')), deleted_at TEXT DEFAULT NULL);
	CREATE TABLE IF NOT EXISTS api_routes (id INTEGER PRIMARY KEY AUTOINCREMENT, route TEXT NOT NULL, method TEXT NOT NULL, description TEXT, status INT2 DEFAULT 1, deleted INT2 DEFAULT 0, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')), deleted_at TEXT DEFAULT NULL, UNIQUE(route, method));
	CREATE TABLE IF NOT EXISTS api_keys (id INTEGER PRIMARY KEY AUTOINCREMENT, key TEXT NOT NULL UNIQUE, description TEXT, status INT2 DEFAULT 1, deleted INT2 DEFAULT 0);
	`
	if _, err := db.Exec(initSQL); err != nil {
		t.Fatalf("failed to create base tables: %v", err)
	}

	if _, err := db.Exec(string(migration)); err != nil {
		t.Fatalf("failed to run migration: %v", err)
	}

	rolesRelQueries := roles_relations.New(db)
	apiKeysQueries := api_keys.New(db)
	apiRoutesQueries := api_routes.New(db)

	repo := repository.NewRolesRelationRepository(rolesRelQueries, apiKeysQueries, apiRoutesQueries)
	registry := security.NewSecurityRegistry(repo)
	if err := registry.Reload(context.Background()); err != nil {
		t.Fatalf("failed to reload security registry: %v", err)
	}

	return db, registry
}

func TestSecurityRegistry_AccessControl(t *testing.T) {
	db, registry := setupTestDB(t)
	defer db.Close()

	// 1. Admin (Role 1) accesses tickets/GetTicketsList/ -> should be ALLOWED
	allowed, enabled := registry.CanAccessRoute([]int64{1}, "POST", "tickets/GetTicketsList/")
	if !enabled || !allowed {
		t.Errorf("expected admin allowed=true, enabled=true, got allowed=%v, enabled=%v", allowed, enabled)
	}

	// 2. Customer (Role 3) accesses tickets/GetTicketsList/ -> should be FORBIDDEN
	allowed, enabled = registry.CanAccessRoute([]int64{3}, "POST", "tickets/GetTicketsList/")
	if !enabled || allowed {
		t.Errorf("expected customer allowed=false on GetTicketsList, got allowed=%v", allowed)
	}

	// 3. Customer (Role 3) accesses tickets/GetTicketByID/ -> should be ALLOWED
	allowed, enabled = registry.CanAccessRoute([]int64{3}, "POST", "tickets/GetTicketByID/")
	if !enabled || !allowed {
		t.Errorf("expected customer allowed=true on GetTicketByID, got allowed=%v", allowed)
	}

	// 4. Unknown route -> allowed by default for authenticated users
	allowed, enabled = registry.CanAccessRoute([]int64{3}, "POST", "unknown/route/")
	if !enabled || !allowed {
		t.Errorf("expected unmapped route allowed=true, got %v", allowed)
	}
}

func TestSecurityRegistry_GetPermissionsForRoles(t *testing.T) {
	db, registry := setupTestDB(t)
	defer db.Close()

	// Admin should have all 9 permissions
	adminPerms := registry.GetPermissionsForRoles([]int64{1})
	if len(adminPerms) != 9 {
		t.Errorf("expected 9 permissions for admin, got %d: %v", len(adminPerms), adminPerms)
	}

	// Customer should have 6 permissions
	custPerms := registry.GetPermissionsForRoles([]int64{3})
	if len(custPerms) != 6 {
		t.Errorf("expected 6 permissions for customer, got %d: %v", len(custPerms), custPerms)
	}
}
