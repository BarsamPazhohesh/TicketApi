package scenario_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/testutil"
	"time"
)

// TestScenario_OTPSendAndVerifyFullLifecycle tests the full end-to-end OTP flow
func TestScenario_OTPSendAndVerifyFullLifecycle(t *testing.T) {
	// 1. Mock SMS Provider Gateway
	smsGatewayCalled := false
	var receivedMessage string
	var receivedUserName string
	var receivedPassWord string
	var receivedSender string
	var receivedReceptor string

	smsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		smsGatewayCalled = true
		_ = r.ParseForm()
		receivedUserName = r.FormValue("userName")
		receivedPassWord = r.FormValue("passWord")
		receivedSender = r.FormValue("senderNumber")
		receivedReceptor = r.FormValue("reciverNumber")
		receivedMessage = r.FormValue("smsText")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("3257258978\n<html>success</html>"))
	}))
	defer smsServer.Close()

	_ = os.Setenv("OTP_URL", smsServer.URL)
	_ = os.Setenv("OTP_USERNAME", "test_user")
	_ = os.Setenv("OTP_PASSWORD", "test_pass")
	_ = os.Setenv("OTP_SENDER_NUMBER", "1000")

	// 2. Setup Test App with SQLite + in-memory state
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	phone := "09121234567"
	captchaToken := testutil.GenerateTestCaptchaTokenWithPhone(t, app.Services.Token, "", "")
	captchaCookieName := config.Get().Captcha.CookieName

	// ─── STEP 1: Client requests OTP via HTTP endpoint ───
	t.Log("🚀 [STEP 1]: Client requests OTP")
	sendBody, _ := json.Marshal(dto.SendOTPDTO{PhoneNumber: phone})
	reqSend, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/send/", bytes.NewReader(sendBody))
	reqSend.Header.Set("Content-Type", "application/json")
	reqSend.AddCookie(&http.Cookie{
		Name:  captchaCookieName,
		Value: captchaToken,
	})
	wSend := httptest.NewRecorder()

	app.Engine.ServeHTTP(wSend, reqSend)
	if wSend.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on send OTP, got %d, body: %s", wSend.Code, wSend.Body.String())
	}

	var sendResp dto.SendOTPResponseDTO
	if err := json.Unmarshal(wSend.Body.Bytes(), &sendResp); err != nil {
		t.Fatalf("failed to decode send OTP response: %v", err)
	}
	if sendResp.Message != "کد تایید با موفقیت ارسال شد" {
		t.Fatalf("expected Persian message 'کد تایید با موفقیت ارسال شد', got: %s", sendResp.Message)
	}

	// ─── STEP 2: Verify async SMS delivery, payload keys, and warehouse record ───
	t.Log("⏳ [STEP 2]: Verify async SMS delivery and warehouse status")
	time.Sleep(150 * time.Millisecond) // wait for async goroutine

	if !smsGatewayCalled {
		t.Fatal("expected SMS gateway to be called")
	}

	if receivedUserName != "test_user" || receivedPassWord != "test_pass" || receivedSender != "1000" || receivedReceptor != phone {
		t.Fatalf("unexpected payload keys/values sent to SMS gateway: user=%s, pass=%s, sender=%s, receptor=%s",
			receivedUserName, receivedPassWord, receivedSender, receivedReceptor)
	}

	record, apiErr := app.Repos.SMSWarehouse.GetByID(context.Background(), 1)
	if apiErr != nil {
		t.Fatalf("failed to query SMS warehouse: %v", apiErr)
	}

	if record.ReceiverPhoneNumber != phone {
		t.Fatalf("expected warehouse phone %s, got %s", phone, record.ReceiverPhoneNumber)
	}

	if record.Status != dto.SMSStatusSent {
		t.Fatalf("expected warehouse status Sent (1), got %d", record.Status)
	}

	// Extract the 6-digit code from the Persian message: "کد تایید شما: XXXXXX"
	parts := strings.Split(receivedMessage, ": ")
	if len(parts) < 2 {
		t.Fatalf("unexpected message format: %s", receivedMessage)
	}
	otpCode := strings.TrimSpace(parts[1])
	t.Logf("🔑 Generated OTP Code: %s", otpCode)

	// ─── STEP 3: Client attempts verification with WRONG code ───
	t.Log("❌ [STEP 3]: Verification with invalid code should fail")
	wrongBody, _ := json.Marshal(dto.VerifyOTPDTO{
		PhoneNumber: phone,
		Code:        "999999",
	})
	reqWrong, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/verify/", bytes.NewReader(wrongBody))
	reqWrong.Header.Set("Content-Type", "application/json")
	reqWrong.AddCookie(&http.Cookie{
		Name:  captchaCookieName,
		Value: captchaToken,
	})
	wWrong := httptest.NewRecorder()

	app.Engine.ServeHTTP(wWrong, reqWrong)
	if wWrong.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for wrong code, got %d", wWrong.Code)
	}

	// ─── STEP 4: Client attempts verification with CORRECT code ───
	t.Log("✅ [STEP 4]: Verification with correct code should succeed")
	correctBody, _ := json.Marshal(dto.VerifyOTPDTO{
		PhoneNumber: phone,
		Code:        otpCode,
	})
	reqCorrect, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/verify/", bytes.NewReader(correctBody))
	reqCorrect.Header.Set("Content-Type", "application/json")
	reqCorrect.AddCookie(&http.Cookie{
		Name:  captchaCookieName,
		Value: captchaToken,
	})
	wCorrect := httptest.NewRecorder()

	app.Engine.ServeHTTP(wCorrect, reqCorrect)
	if wCorrect.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for correct code, got %d, body: %s", wCorrect.Code, wCorrect.Body.String())
	}

	var verifyResp dto.VerifyOTPResponseDTO
	if err := json.Unmarshal(wCorrect.Body.Bytes(), &verifyResp); err != nil {
		t.Fatalf("failed to decode verify response: %v", err)
	}
	if !verifyResp.Valid {
		t.Fatal("expected valid=true in verification response")
	}
	if verifyResp.Message != "تایید با موفقیت انجام شد" {
		t.Fatalf("expected Persian message 'تایید با موفقیت انجام شد', got: %s", verifyResp.Message)
	}

	// Verify Level 2 elevated cookie
	cookies := wCorrect.Result().Cookies()
	var elevatedCaptchaToken string
	for _, ck := range cookies {
		if ck.Name == captchaCookieName {
			elevatedCaptchaToken = ck.Value
			break
		}
	}
	if elevatedCaptchaToken == "" {
		t.Fatal("expected captcha_token cookie with verified phone on verify success")
	}
	claims, errParse := app.Services.Token.ParseCaptchaToken(elevatedCaptchaToken)
	if errParse != nil || claims.PhoneNumber != phone {
		t.Fatalf("expected elevated captcha token with phone %s, got: %+v", phone, claims)
	}

	// ─── STEP 5: Replay attack check (code must be single-use) ───
	t.Log("🔒 [STEP 5]: Reusing consumed code should fail")
	reqReplay, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/verify/", bytes.NewReader(correctBody))
	reqReplay.Header.Set("Content-Type", "application/json")
	reqReplay.AddCookie(&http.Cookie{
		Name:  captchaCookieName,
		Value: captchaToken,
	})
	wReplay := httptest.NewRecorder()

	app.Engine.ServeHTTP(wReplay, reqReplay)
	if wReplay.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for already consumed code, got %d", wReplay.Code)
	}

	t.Log("🎉 [SCENARIO PASSED]: Full OTP send, validation, and single-use lifecycle verified.")
}

// TestScenario_OTPMaxVerificationAttemptsLockout tests the 5-attempt lockout and OTP code eviction
func TestScenario_OTPMaxVerificationAttemptsLockout(t *testing.T) {
	// 1. Mock SMS Provider Gateway
	smsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("123456\n<html>success</html>"))
	}))
	defer smsServer.Close()

	_ = os.Setenv("OTP_URL", smsServer.URL)
	_ = os.Setenv("OTP_USERNAME", "test_user")
	_ = os.Setenv("OTP_PASSWORD", "test_pass")
	_ = os.Setenv("OTP_SENDER_NUMBER", "1000")

	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	phone := "09127778899"
	captchaToken := testutil.GenerateTestCaptchaTokenWithPhone(t, app.Services.Token, "", "")
	captchaCookieName := config.Get().Captcha.CookieName

	// ─── Step 1: Send OTP ───
	sendBody, _ := json.Marshal(dto.SendOTPDTO{PhoneNumber: phone})
	reqSend, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/send/", bytes.NewReader(sendBody))
	reqSend.Header.Set("Content-Type", "application/json")
	reqSend.AddCookie(&http.Cookie{Name: captchaCookieName, Value: captchaToken})
	wSend := httptest.NewRecorder()
	app.Engine.ServeHTTP(wSend, reqSend)
	if wSend.Code != http.StatusOK {
		t.Fatalf("expected 200 on SendOTP, got %d", wSend.Code)
	}

	wrongVerifyBody, _ := json.Marshal(dto.VerifyOTPDTO{
		PhoneNumber: phone,
		Code:        "000000",
	})

	// ─── Step 2: Attempts 1 through 4 should fail with ErrOTPInvalid (400) ───
	for attempt := 1; attempt <= 4; attempt++ {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/verify/", bytes.NewReader(wrongVerifyBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: captchaCookieName, Value: captchaToken})
		w := httptest.NewRecorder()
		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("attempt %d: expected 400 Bad Request, got %d", attempt, w.Code)
		}
	}

	// ─── Step 3: Attempt 5 should reach lockout limit and return ErrOTPMaxAttemptsExceeded (400) ───
	req5, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/verify/", bytes.NewReader(wrongVerifyBody))
	req5.Header.Set("Content-Type", "application/json")
	req5.AddCookie(&http.Cookie{Name: captchaCookieName, Value: captchaToken})
	w5 := httptest.NewRecorder()
	app.Engine.ServeHTTP(w5, req5)
	if w5.Code != http.StatusBadRequest {
		t.Fatalf("attempt 5: expected 400 Bad Request, got %d", w5.Code)
	}
	if !strings.Contains(w5.Body.String(), "حد مجاز خطا در وارد کردن کد یکبار مصرف به پایان رسید") {
		t.Fatalf("attempt 5: expected Persian max attempts error message, got: %s", w5.Body.String())
	}

	// ─── Step 4: Attempt 6 should immediately return ErrOTPExpired (code was purged from cache) ───
	req6, _ := http.NewRequest(http.MethodPost, "/api/v1/otp/verify/", bytes.NewReader(wrongVerifyBody))
	req6.Header.Set("Content-Type", "application/json")
	req6.AddCookie(&http.Cookie{Name: captchaCookieName, Value: captchaToken})
	w6 := httptest.NewRecorder()
	app.Engine.ServeHTTP(w6, req6)
	if w6.Code != http.StatusBadRequest {
		t.Fatalf("attempt 6: expected 400 Bad Request, got %d", w6.Code)
	}
	if !strings.Contains(w6.Body.String(), "کد تایید منقضی شده است") {
		t.Fatalf("attempt 6: expected ErrOTPExpired message after code eviction, got: %s", w6.Body.String())
	}
}
