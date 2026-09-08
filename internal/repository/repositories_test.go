package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/db/ticket_priorities"
	"ticket-api/internal/db/ticket_statuses"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/dto"
	"ticket-api/internal/testutil"
)

func TestSQLRepositories_CRUDAndIntegrity(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	ctx := context.Background()

	// 1. Users Repository
	t.Run("Users Repository", func(t *testing.T) {
		res, err := app.Repos.Users.CreateUserWithPassword(ctx, dto.SignUpWithPasswordDTO{
			Username:     "09121112233",
			Password:     "hashpass1234",
			DepartmentID: 1,
		})
		if err != nil {
			t.Fatalf("unexpected error creating user: %v", err)
		}
		if res == nil || res.ID <= 0 {
			t.Fatalf("expected positive user ID, got %v", res)
		}

		userID := res.ID
		exists, err := app.Repos.Users.IsUserExist(ctx, userID)
		if err != nil || !exists {
			t.Fatalf("expected user %d to exist, err: %v", userID, err)
		}

		user, err := app.Repos.Users.GetUserByID(ctx, userID)
		if err != nil || user == nil || user.Username != "09121112233" {
			t.Fatalf("expected valid user, got %v, err %v", user, err)
		}

		byUser, err := app.Repos.Users.GetUserByUsername(ctx, "09121112233")
		if err != nil || byUser == nil || byUser.ID != userID {
			t.Fatalf("expected user by username, got %v, err %v", byUser, err)
		}
	})

	// 2. Departments Repository with Cache
	t.Run("Departments Repository", func(t *testing.T) {
		depts, err := app.Repos.Departments.GetAllDepartments(ctx)
		if err != nil {
			t.Fatalf("failed getting departments: %v", err)
		}
		if len(depts) == 0 {
			t.Fatalf("expected seeded departments, got 0")
		}

		newID, err := app.Repos.Departments.AddDepartment(ctx, departments.AddDepartmentParams{
			Title:       "Finance",
			Description: sql.NullString{String: "Finance dept", Valid: true},
		})
		if err != nil || newID <= 0 {
			t.Fatalf("failed adding department: %v", err)
		}

		exists, err := app.Repos.Departments.IsDepartmentExits(ctx, newID)
		if err != nil || !exists {
			t.Fatalf("expected new department %d to exist", newID)
		}
	})

	// 3. Ticket Types Repository
	t.Run("Ticket Types Repository", func(t *testing.T) {
		types, err := app.Repos.TicketTypes.GetAllTicketTypes(ctx)
		if err != nil || len(types) == 0 {
			t.Fatalf("expected ticket types, got %v, err: %v", types, err)
		}

		activeTypes, err := app.Repos.TicketTypes.GetAllActiveTicketTypes(ctx)
		if err != nil || len(activeTypes) == 0 {
			t.Fatalf("expected active ticket types, got %v, err: %v", activeTypes, err)
		}

		newTypeID, err := app.Repos.TicketTypes.AddTicketType(ctx, ticket_types.AddTicketTypeParams{
			Title:       "Security Incident",
			Description: sql.NullString{String: "Security issues", Valid: true},
		})
		if err != nil || newTypeID <= 0 {
			t.Fatalf("failed adding ticket type: %v", err)
		}

		exists, err := app.Repos.TicketTypes.IsTicketTypeExits(ctx, newTypeID)
		if err != nil || !exists {
			t.Fatalf("expected ticket type %d to exist", newTypeID)
		}
	})

	// 4. Ticket Statuses Repository
	t.Run("Ticket Statuses Repository", func(t *testing.T) {
		statuses, err := app.Repos.TicketStatus.GetAllActiveTicketStatuses(ctx)
		if err != nil || len(statuses) == 0 {
			t.Fatalf("expected ticket statuses, got %v, err: %v", statuses, err)
		}

		openStatus, err := app.Repos.TicketStatus.GetOpenStatus(ctx)
		if err != nil || openStatus == nil || openStatus.ID <= 0 {
			t.Fatalf("expected open status, got %v, err: %v", openStatus, err)
		}

		addErr := app.Repos.TicketStatus.AddTicketStatus(ctx, ticket_statuses.AddTicketStatusParams{
			Title:       "escalated",
			Description: sql.NullString{String: "Escalated to level 2", Valid: true},
		})
		if addErr != nil {
			t.Fatalf("failed adding ticket status: %v", addErr)
		}
	})

	// 5. Ticket Priorities Repository
	t.Run("Ticket Priorities Repository", func(t *testing.T) {
		addErr := app.Repos.TicketPriorities.AddTicketPriority(ctx, ticket_priorities.AddTicketPriorityParams{
			UserID:       2,
			TicketTypeID: 2,
			Priority:     2,
		})
		if addErr != nil {
			t.Fatalf("failed adding ticket priority: %v", addErr)
		}
	})

	// 6. Roles & RolesRelations Repository
	t.Run("Roles and Relations", func(t *testing.T) {
		roleID, err := app.Repos.Roles.AddRole(ctx, "qa_engineer")
		if err != nil || roleID <= 0 {
			t.Fatalf("failed adding role: %v", err)
		}

		exists, err := app.Repos.Roles.IsRoleExist(ctx, roleID)
		if err != nil || !exists {
			t.Fatalf("expected role %d to exist", roleID)
		}

		rolesWithPerms, err := app.Repos.RolesRelations.GetAllRolesWithPermissions(ctx)
		if err != nil {
			t.Fatalf("failed getting roles with perms: %v", err)
		}
		_ = rolesWithPerms

		routesWithPerms, err := app.Repos.RolesRelations.GetAllRoutesWithPermissions(ctx)
		if err != nil {
			t.Fatalf("failed getting routes with perms: %v", err)
		}
		_ = routesWithPerms
	})
}
