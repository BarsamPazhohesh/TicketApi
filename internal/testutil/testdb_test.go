package testutil_test

import (
	"testing"
	"ticket-api/internal/testutil"
)

func TestSetupTestApp(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	if app.DB == nil {
		t.Fatal("expected non-nil DB")
	}
	if app.Engine == nil {
		t.Fatal("expected non-nil Gin engine")
	}

	token := testutil.GenerateTestAuthToken(t, app.Services.Token, 1, "admin", []int64{1})
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}
