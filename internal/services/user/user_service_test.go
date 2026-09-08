package user_test

import (
	"context"
	"testing"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/testutil"
)

func TestUserService_GetUsers(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	ctx := context.Background()

	// Seed test users
	u1, err := app.Repos.Users.CreateUserWithPassword(ctx, dto.SignUpWithPasswordDTO{
		Username:     "09121110001",
		Password:     "pass112345",
		DepartmentID: 1,
	})
	if err != nil {
		t.Fatalf("failed creating user 1: %v", err)
	}

	u2, err := app.Repos.Users.CreateUserWithPassword(ctx, dto.SignUpWithPasswordDTO{
		Username:     "09121110002",
		Password:     "pass212345",
		DepartmentID: 1,
	})
	if err != nil {
		t.Fatalf("failed creating user 2: %v", err)
	}

	// 1. GetUserByID
	t.Run("GetUserByID valid & invalid", func(t *testing.T) {
		dtoUser, apiErr := app.Services.User.GetUserByID(ctx, u1.ID)
		if apiErr != nil || dtoUser == nil {
			t.Fatalf("expected user, got %v, err %v", dtoUser, apiErr)
		}
		if dtoUser.Username != "09121110001" {
			t.Errorf("expected username 09121110001, got %s", dtoUser.Username)
		}

		_, notFoundErr := app.Services.User.GetUserByID(ctx, 99999)
		if notFoundErr == nil {
			t.Fatal("expected error on non-existent user")
		}
	})

	// 2. GetUsersByIDs
	t.Run("GetUsersByIDs", func(t *testing.T) {
		usersList, apiErr := app.Services.User.GetUsersByIDs(ctx, []int64{u1.ID, u2.ID})
		if apiErr != nil {
			t.Fatalf("unexpected error: %v", apiErr)
		}
		if len(usersList) != 2 {
			t.Fatalf("expected 2 users, got %d", len(usersList))
		}
	})

	// 3. GetUserByUsername
	t.Run("GetUserByUsername", func(t *testing.T) {
		dtoUser, apiErr := app.Services.User.GetUserByUsername(ctx, "09121110002")
		if apiErr != nil || dtoUser == nil {
			t.Fatalf("unexpected error: %v", apiErr)
		}
		if dtoUser.ID != u2.ID {
			t.Errorf("expected id %d, got %d", u2.ID, dtoUser.ID)
		}

		_, apiErr = app.Services.User.GetUserByUsername(ctx, "nonexistent")
		if apiErr == nil || apiErr.Err.Code != errx.ErrUserNotFound {
			t.Fatalf("expected user not found error, got %v", apiErr)
		}
	})
}
