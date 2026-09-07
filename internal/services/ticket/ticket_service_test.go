package ticket_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"ticket-api/internal/config"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/db/ticket_priorities"
	"ticket-api/internal/db/ticket_statuses"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/db/users"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/repository"
	"ticket-api/internal/services/cache"
	"ticket-api/internal/services/ticket"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	_ = os.Setenv("JWT_SECRET", "this-is-a-very-secure-32-byte-long-secret-key-12345")
	errx.NewRegistry(nil)
}

func setupTestDB(t *testing.T) (*sql.DB, *repository.UsersRepository, *repository.TicketTypesRepository, *repository.DepartmentsRepository, *repository.TicketStatusesRepository, *repository.TicketPrioritiesRepository) {
	config.Load("../../../config.yaml")

	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory: %v", err)
	}

	schemas := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			deleted_at TEXT DEFAULT NULL,
			status INTEGER NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE ticket_types (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			deleted_at TEXT DEFAULT NULL,
			status INTEGER NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE departments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			deleted_at TEXT DEFAULT NULL,
			status INTEGER NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE ticket_statuses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			deleted_at TEXT DEFAULT NULL,
			status INTEGER NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE ticket_priorities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			deleted_at TEXT DEFAULT NULL,
			status INTEGER NOT NULL DEFAULT 1
		);`,
	}

	for _, s := range schemas {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("failed schema execution: %v", err)
		}
	}

	cacheSvc := cache.NewCacheService(nil)

	userRepo := repository.NewUsersRepository(users.New(db))
	ticketTypeRepo := repository.NewTicketTypesRepository(ticket_types.New(db), cacheSvc)
	departmentRepo := repository.NewDepartmentsRepository(departments.New(db), cacheSvc)
	ticketStatusRepo := repository.NewTicketStatusesRepository(ticket_statuses.New(db), cacheSvc)
	ticketPriorityRepo := repository.NewTicketPrioritiesRepository(ticket_priorities.New(db))

	return db, userRepo, ticketTypeRepo, departmentRepo, ticketStatusRepo, ticketPriorityRepo
}

func TestTicketService_ValidationLimits(t *testing.T) {
	db, userRepo, ticketTypeRepo, departmentRepo, ticketStatusRepo, ticketPriorityRepo := setupTestDB(t)
	defer db.Close()

	svc := ticket.NewTicketService(nil, nil, ticketTypeRepo, ticketPriorityRepo, ticketStatusRepo, userRepo, departmentRepo, nil)

	ctx := context.Background()

	// 1. Exceeded attachments limit
	req := dto.TicketCreateRequest{
		TicketTypeID: 1,
		DepartmentID: 1,
		Title:        "Test Title",
		Body:         "Test Body",
		Attachments:  []string{"http://storage/temp/1.png", "http://storage/temp/2.png", "http://storage/temp/3.png", "http://storage/temp/4.png", "http://storage/temp/5.png", "http://storage/temp/6.png"},
	}

	_, apiErr := svc.CreateTicket(ctx, 1, "", req)
	if apiErr == nil {
		t.Fatalf("expected error for exceeding attachment limit, got nil")
	}
	if apiErr.Err.Code != errx.ErrMaxTicketFilesExceeded {
		t.Fatalf("expected code %d, got %d", errx.ErrMaxTicketFilesExceeded, apiErr.Err.Code)
	}

	// 2. User not found
	req.Attachments = nil
	_, apiErr = svc.CreateTicket(ctx, 999, "", req)
	if apiErr == nil {
		t.Fatalf("expected error for non-existent user, got nil")
	}
	if apiErr.Err.Code != errx.ErrUserNotFound {
		t.Fatalf("expected code %d, got %d", errx.ErrUserNotFound, apiErr.Err.Code)
	}

	// 3. Ticket type not found
	_, err := db.Exec(`INSERT INTO users (id, username, password) VALUES (1, 'john_doe', 'hash');`)
	if err != nil {
		t.Fatalf("failed insert user: %v", err)
	}

	_, apiErr = svc.CreateTicket(ctx, 1, "", req)
	if apiErr == nil {
		t.Fatalf("expected error for non-existent ticket type, got nil")
	}
	if apiErr.Err.Code != errx.ErrTicketTypeNotFound {
		t.Fatalf("expected code %d, got %d", errx.ErrTicketTypeNotFound, apiErr.Err.Code)
	}

	// 4. Department not found
	_, err = db.Exec(`INSERT INTO ticket_types (id, title) VALUES (1, 'Technical');`)
	if err != nil {
		t.Fatalf("failed insert ticket type: %v", err)
	}

	_, apiErr = svc.CreateTicket(ctx, 1, "", req)
	if apiErr == nil {
		t.Fatalf("expected error for non-existent department, got nil")
	}
	if apiErr.Err.Code != errx.ErrDepartmentNotFound {
		t.Fatalf("expected code %d, got %d", errx.ErrDepartmentNotFound, apiErr.Err.Code)
	}
}

func TestTicketService_GetTicketByTrackCode_InvalidFormat(t *testing.T) {
	db, userRepo, ticketTypeRepo, departmentRepo, ticketStatusRepo, ticketPriorityRepo := setupTestDB(t)
	defer db.Close()

	svc := ticket.NewTicketService(nil, nil, ticketTypeRepo, ticketPriorityRepo, ticketStatusRepo, userRepo, departmentRepo, nil)
	ctx := context.Background()

	_, apiErr := svc.GetTicketByTrackCode(ctx, dto.TicketByTrackCodeRequestDTO{
		TrackCode: "invalid-track-code",
	}, 1, "", nil)
	if apiErr == nil {
		t.Fatalf("expected error on invalid track code format")
	}
	if apiErr.Err.Code != errx.ErrBadRequest {
		t.Fatalf("expected ErrBadRequest (%d), got %d", errx.ErrBadRequest, apiErr.Err.Code)
	}
}

func TestTicketService_GetActiveTypesAndStatuses(t *testing.T) {
	db, userRepo, ticketTypeRepo, departmentRepo, ticketStatusRepo, ticketPriorityRepo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	_, _ = db.Exec(`INSERT INTO ticket_types (id, title, status) VALUES (1, 'Billing', 1), (2, 'Technical', 0);`)
	_, _ = db.Exec(`INSERT INTO ticket_statuses (id, title, status) VALUES (1, 'Open', 1), (2, 'Closed', 1);`)

	svc := ticket.NewTicketService(nil, nil, ticketTypeRepo, ticketPriorityRepo, ticketStatusRepo, userRepo, departmentRepo, nil)

	types, apiErr := svc.GetAllActiveTicketTypes(ctx)
	if apiErr != nil {
		t.Fatalf("unexpected error getting ticket types: %v", apiErr)
	}
	if len(types) != 1 {
		t.Fatalf("expected 1 active ticket type, got %d", len(types))
	}
	if types[0].Title != "Billing" {
		t.Fatalf("expected Billing, got %s", types[0].Title)
	}

	statuses, apiErr := svc.GetAllActiveTicketStatuses(ctx)
	if apiErr != nil {
		t.Fatalf("unexpected error getting ticket statuses: %v", apiErr)
	}
	if len(statuses) != 2 {
		t.Fatalf("expected 2 active ticket statuses, got %d", len(statuses))
	}
}
