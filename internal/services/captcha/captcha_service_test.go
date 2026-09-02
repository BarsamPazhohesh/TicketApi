package captcha_test

import (
	"context"
	"testing"
	"ticket-api/internal/testutil"
)

func TestCaptchaService_GenerateAndVerify(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	ctx := context.Background()

	// 1. Generate captcha ID
	res, err := app.Services.Captcha.GenerateCaptcha()
	if err != nil || res == nil || res.ID == "" {
		t.Fatalf("expected valid captcha result, got %v, err: %v", res, err)
	}

	// 2. Verify invalid digits fail
	apiErr := app.Services.Captcha.VerifyCaptcha(res.ID, "000000")
	if apiErr == nil {
		t.Fatal("expected captcha verification to fail for invalid solution")
	}

	// 3. Captcha Token generation and parsing
	ip := "127.0.0.1"
	tokenStr, apiErr := app.Services.Token.NewCaptchaToken(ip)
	if apiErr != nil || tokenStr == "" {
		t.Fatalf("failed generating captcha token: %v", apiErr)
	}

	claims, apiErr := app.Services.Token.ParseCaptchaToken(tokenStr)
	if apiErr != nil || claims == nil {
		t.Fatalf("failed parsing captcha token: %v", apiErr)
	}
	if claims.IP != ip {
		t.Fatalf("expected IP %s, got %s", ip, claims.IP)
	}

	// 4. Invalid token parsing
	_, invalidErr := app.Services.Token.ParseCaptchaToken("invalid.jwt.token")
	if invalidErr == nil {
		t.Fatal("expected error for invalid captcha token")
	}
	_ = ctx
}
