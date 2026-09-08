package token_test

import (
	"context"
	"os"
	"testing"
	"ticket-api/internal/config"
	"ticket-api/internal/errx"
	"ticket-api/internal/services/token"
	"time"
)

func init() {
	_ = os.Setenv("JWT_SECRET", "this-is-a-very-secure-32-byte-long-secret-key-12345")
	errx.NewRegistry(nil)
}

func TestAuthToken_GenerateAndParse(t *testing.T) {
	config.Load("../../../config.yaml")

	tokenService := token.NewTokenService(nil)
	ctx := context.Background()

	claims := token.AuthClaims{
		UserID:      42,
		Username:    "testuser",
		PhoneNumber: "09123456789",
		RoleIDs:     []int64{1, 2},
		Permissions: []string{"tickets:read", "tickets:create"},
	}

	tokenStr, apiErr := tokenService.NewAuthToken(claims)
	if apiErr != nil {
		t.Fatalf("expected nil error, got %v", apiErr)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token string")
	}

	parsed, parseErr := tokenService.ParseAuthToken(ctx, tokenStr)
	if parseErr != nil {
		t.Fatalf("expected nil parse error, got %v", parseErr)
	}
	if parsed == nil {
		t.Fatal("expected parsed claims not nil")
	}
	if parsed.UserID != 42 || parsed.Username != "testuser" || parsed.PhoneNumber != "09123456789" {
		t.Errorf("expected UserID=42, Username='testuser', PhoneNumber='09123456789', got %d, %s, %s", parsed.UserID, parsed.Username, parsed.PhoneNumber)
	}
	if len(parsed.RoleIDs) != 2 || parsed.RoleIDs[0] != 1 || parsed.RoleIDs[1] != 2 {
		t.Errorf("unexpected RoleIDs: %v", parsed.RoleIDs)
	}
	if parsed.ID == "" {
		t.Error("expected non-empty JTI claim")
	}
}

func TestAuthToken_InvalidAndTamperedToken(t *testing.T) {
	config.Load("../../../config.yaml")
	tokenService := token.NewTokenService(nil)
	ctx := context.Background()

	// 1. Random string
	parsed, err := tokenService.ParseAuthToken(ctx, "invalid.jwt.token")
	if err == nil || parsed != nil {
		t.Errorf("expected error for invalid token, got parsed=%v", parsed)
	}

	// 2. Empty string
	parsed, err = tokenService.ParseAuthToken(ctx, "")
	if err == nil || parsed != nil {
		t.Errorf("expected error for empty token string")
	}
}

func TestAuthToken_Revocation(t *testing.T) {
	config.Load("../../../config.yaml")
	tokenService := token.NewTokenService(nil)
	ctx := context.Background()

	claims := token.AuthClaims{
		UserID:   99,
		Username: "revokeme",
		RoleIDs:  []int64{3},
	}

	tokenStr, _ := tokenService.NewAuthToken(claims)
	parsed, err := tokenService.ParseAuthToken(ctx, tokenStr)
	if err != nil || parsed == nil {
		t.Fatalf("failed to parse valid token: %v", err)
	}

	// nil redis safe fallback
	errRevoke := tokenService.RevokeToken(ctx, parsed.ID, 5*time.Minute)
	if errRevoke != nil {
		t.Errorf("expected nil error on nil redis revoke, got %v", errRevoke)
	}
	isRevoked, errCheck := tokenService.IsTokenRevoked(ctx, parsed.ID)
	if errCheck != nil {
		t.Errorf("unexpected error on nil redis check: %v", errCheck)
	}
	if isRevoked {
		t.Error("expected false for IsTokenRevoked with nil redis")
	}
}
