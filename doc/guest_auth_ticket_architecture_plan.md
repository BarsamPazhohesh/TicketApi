# Plan: Ticket API Architecture & HTMX Demo Panel

## 1. Context & Goals
1. **Sender ID Decoupling & Role Derivation**:
   - Strip client-supplied `senderId` from `ChatMessageCreateRequest`.
   - Server derives sender identity (`senderId`: 0 for guest, `user.id` for auth; `senderType`: `"user"` | `"agent"`) via JWT or guest session.
2. **Guest vs Authenticated Handling**:
   - Tickets in MongoDB store `userId` (`0` for guest) and `phoneNumber` at root level for high-performance indexing.
   - Guest ticket lookup verified via `trackCode` + `phoneNumber` + `captcha`.
3. **Demo HTMX Panel (`demo/`)**:
   - Simple, clean web UI built with Go `html/template` + HTMX + Tailwind CSS (via CDN).
   - Reusable template components (buttons, inputs, alerts, modals, chat bubbles, stat cards).
   - Documentation cataloging all UI components for easy extension.
   - **4 Core Views**:
     1. **Guest Dashboard**: Submit ticket with captcha, lookup ticket via track code + phone + captcha, live chat.
     2. **Auth User Dashboard**: Login/register simulation, create ticket, view my tickets, live chat.
     3. **Admin / Agent Dashboard**: Ticket triage, change status, assign department, reply as agent.
     4. **Data Center (Debug)**: Live status & metrics of SQLite, MongoDB, Redis, and MinIO.
   - **Environment Toggle**: Automatically disabled when `GIN_MODE=release` or `ENABLE_DEMO=false` in `.env`.

---

## 2. Reusable Component Catalog (`demo/components/`)
Create DRY sub-templates:
- `components/button.html`: Primary, secondary, danger, ghost buttons with HTMX spinners.
- `components/input.html`: Standard form inputs with labels, errors, and validation classes.
- `components/modal.html`: HTMX-driven dynamic dialogs/modals.
- `components/alert.html`: Success, warning, error toast/banner messages.
- `components/badge.html`: Role (`guest`, `user`, `agent`) and ticket status badges.
- `components/chat_bubble.html`: Left (incoming/agent) vs right (my message) layout with attachments.
- `components/stat_card.html`: Metrics & DB status indicators for Data Center.

---

## 3. Demo Documentation (`demo/README.md`)
Documentation listing:
- Component inventory, parameters, and usage patterns.
- How HTMX endpoints interact with backend API handlers.
- Configuration and security toggles (`ENABLE_DEMO`, `GIN_MODE`).

---

## 4. Implementation Steps

### Step 1: Core API & Schema Adjustments
1. Decouple `ChatMessageCreateRequest`: remove `senderId`.
2. Update `model.Ticket` and `model.ChatMessage` in MongoDB:
   - `Ticket.PhoneNumber string`
   - `ChatMessage.SenderType string` (`"user"` | `"agent"`)
3. Update `TicketService`:
   - `CreateTicket`: Support guest creation (`userID = 0`, `phoneNumber`).
   - `CreateChatMessage`: Automatically assign `senderID` and `senderType`.
   - `GetGuestTicketByTrackCode`: Verify `trackCode` + `phoneNumber` + `captcha`.

### Step 2: HTMX Component Library & Documentation
1. Create `demo/templates/layout.html` (base page with HTMX, Tailwind, dark/light styling).
2. Create `demo/templates/components/*.html` (buttons, inputs, modals, chat bubbles, stat cards).
3. Create `demo/README.md` (component catalog and usage guide).

### Step 3: Demo Views & Handlers
1. Create `demo/templates/views/`:
   - `guest.html` (Submit ticket, track ticket, chat interface).
   - `user.html` (Auth user ticket management).
   - `admin.html` (Staff ticket list, status changer, agent chat).
   - `datacenter.html` (SQLite stats, Mongo stats, Redis ping, MinIO status).
2. Create `internal/handler/demo_handler.go` and routes in `cmd/api/main.go` / `internal/routes/`:
   - Gated behind `os.Getenv("GIN_MODE") != "release" && env.GetEnvBool("ENABLE_DEMO", true)`.

### Step 4: Verification & Docs Sync
1. Run all unit tests: `go test ./...`.
2. Sync Swagger specs: `swag init -g cmd/api/main.go -o docs`.
3. Verify demo interface in browser.
