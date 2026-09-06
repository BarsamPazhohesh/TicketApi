# Guest Ticket & Step-Up OTP Authentication Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement two-level step-up authentication for guest tickets: Level 1 (Captcha token without phone) protects OTP sending against SMS spam, and Level 2 (Captcha token elevated with verified phone on OTP verification) is required for guest ticket creation, tracking, and chat.

**Architecture:** 
1. `VerifyCaptcha` validates captcha and issues Level 1 cookie/token with `IP` and empty `phone_number`.
2. `POST /otp/send/` and `POST /otp/verify/` require Level 1 Captcha token (or Auth token).
3. `POST /otp/verify/` validates code and issues/updates Level 2 Captcha token cookie containing the verified `phone_number`.
4. `CaptchaMiddleware` enforces `RequirePhone: true` for guest ticket routes, ensuring guests must have verified phone in token claims (or be logged in via Auth token).
5. Handlers/Services enforce that guest actions bind strictly to the token's verified `PhoneNumber`.

**Tech Stack:** Go 1.24, Gin, JWT (golang-jwt/v5), Redis cache, MinIO, SQLite, swag.

---

### Task 1: Update Captcha DTO and VerifyCaptcha Handler (Level 1 Token)

**Files:**
- Modify: `internal/dto/captcha_dto.go:3-7`
- Modify: `internal/handler/captcha_handler.go:57-78`
- Test: `internal/handler/captcha_handler_test.go`

- [ ] **Step 1: Update `dto.CaptchaVerifyRequest` to remove `PhoneNumber`**
- [ ] **Step 2: Update `VerifyCaptchaHandler` to issue Level 1 token with empty phone (`""`)**
- [ ] **Step 3: Update `captcha_handler_test.go` and run `go test ./internal/handler/... -run TestCaptchaHandler`**
- [ ] **Step 4: Commit changes**

---

### Task 2: Update OTP Handler to require Level 1 token and issue Level 2 token on verify

**Files:**
- Modify: `internal/handler/otp_handler.go`
- Modify: `internal/handler/app_handlers.go:19,31`
- Test: `internal/handler/otp_handler_test.go`

- [ ] **Step 1: Pass `TokenService` to `OTPHandler`**
- [ ] **Step 2: On successful `VerifyOTP`, generate Level 2 token (`phone = req.PhoneNumber`) and set cookie**
- [ ] **Step 3: Update `otp_handler_test.go` to test cookie issuance on verify**
- [ ] **Step 4: Commit changes**

---

### Task 3: Enhance `CaptchaMiddleware` to support Level 1 vs Level 2 (Required Phone)

**Files:**
- Modify: `internal/middleware/captcha_middleware.go`
- Modify: `cmd/api/routes.go`
- Test: `internal/middleware/access_control_test.go`

- [ ] **Step 1: Add options / parameter to `CaptchaMiddleware(tokenService, requirePhone bool)`**
- [ ] **Step 2: Attach Level 1 guard to OTP routes (`/otp/send`, `/otp/verify`)**
- [ ] **Step 3: Attach Level 2 guard (`requirePhone: true`) to guest ticket routes (`/tickets/CreateTicket/`, `/tickets/GetTicketByTrackCode/`, `/tickets/:id/CreateChat/`)**
- [ ] **Step 4: Pass verified phone to Gin context (`c.Set("guest_phone", claims.PhoneNumber)`)**
- [ ] **Step 5: Commit changes**

---

### Task 4: Bind Guest Ticket and Chat Actions to Verified Phone

**Files:**
- Modify: `internal/handler/ticket_handler.go`
- Modify: `internal/handler/chat_handlers.go`
- Test: `internal/handler/ticket_handler_test.go`

- [ ] **Step 1: In `CreateTicketHandler`, if guest (`currentUserID == 0`), override `req.PhoneNumber` from `guest_phone` context**
- [ ] **Step 2: In `GetTicketByTrackCodeHandler`, if guest, override `req.PhoneNumber` from `guest_phone` context**
- [ ] **Step 3: In `CreateChatHandler`, enforce guest ownership if `senderID == 0`**
- [ ] **Step 4: Commit changes**

---

### Task 5: End-to-End Scenario & Regression Tests and Swagger Update

**Files:**
- Modify: `test/scenario/ticket_flow_test.go`
- Modify: `test/scenario/otp_scenario_test.go`
- Docs: `docs/` (via `swag init -g cmd/api/main.go -o docs`)

- [ ] **Step 1: Update end-to-end scenario tests to follow Captcha -> OTP Send -> OTP Verify -> Ticket Create flow**
- [ ] **Step 2: Run all unit and scenario tests (`go test ./...`)**
- [ ] **Step 3: Regenerate Swagger docs (`swag init -g cmd/api/main.go -o docs`)**
- [ ] **Step 4: Commit changes**
