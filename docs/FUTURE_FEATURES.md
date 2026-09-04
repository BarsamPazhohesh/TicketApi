# Future Roadmap & Backlog

This document tracks planned architectural discussions, improvements, and test coverage expansions for upcoming iterations.

---

## 1. Background Worker Architecture & Distributed Processing
- [ ] **Worker Pattern Review**:
  - Current state: In-process goroutine with `time.NewTicker` (`startPendingSMSRetryWorker`).
  - Discussion points:
    - Distributed locking (Redis `SETNX` / redlock) to prevent duplicate SMS dispatch when running multiple backend API instances in production.
    - Max retry count & backoff strategy (exponential backoff vs fixed 5-min intervals) to avoid endless retries on invalid recipient numbers.
    - Dead Letter Queue (DLQ) or `status = 3 (abandoned)` for un-deliverable SMS after $N$ attempts.
    - Rate limiting per SMS provider gateway constraints.

---

## 2. Test Coverage & Integrity Expansions
- [ ] **Higher Unit & Integration Coverage**:
  - Target: Push `internal/services/otp` coverage from ~58% to >85%.
  - Unit test `retryPendingSMS` worker loop directly (mocking failures, retries, and status transitions to `SMSStatusFailed` and `SMSStatusSent`).
  - Unit test `SMSLoaderService` (`loadFromDB`, `autoRefresh`, fallback handling when DB template is missing).
  - Test network timeout and HTTP 5xx errors from the SMS provider.

---

## 3. End-to-End Scenario Testing Expansion
- [ ] **Multi-step User Journeys**:
  - OTP Sign-up / Login journey: `POST /otp/send` -> `POST /otp/verify` -> issue Auth JWT -> access protected routes.
  - OTP Rate limiting & spam protection: assert max requests per IP/phone number within a sliding window.
  - Gateway failover scenario: test pending warehouse queue draining after mock SMS gateway recovery.
  - SMS template hot-reload test: change `sms_messages` in DB and verify `SMSLoaderService` reflects updated text without restarting the server.
