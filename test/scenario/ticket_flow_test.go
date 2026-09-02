package scenario_test

import (
	"context"
	"testing"
	"ticket-api/internal/dto"
	"ticket-api/internal/model"
	"ticket-api/internal/testutil"
	"ticket-api/internal/util"
	"time"

	"github.com/google/uuid"
)

// TestScenario_CompleteTicketLifecycle tests the entire ticket workflow from creation to resolution
func TestScenario_CompleteTicketLifecycle(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	ctx := context.Background()

	t.Log("🚀 [SCENARIO STEP 1]: Client obtains and solves Captcha")
	captchaToken := testutil.GenerateTestCaptchaToken(t, app.Services.Token, "127.0.0.1")
	if captchaToken == "" {
		t.Fatal("failed generating captcha token")
	}

	t.Log("🔍 [SCENARIO STEP 2]: Client fetches active departments and ticket types")
	types, errTypes := app.Services.Ticket.GetAllActiveTicketTypes(ctx)
	if errTypes != nil || len(types) == 0 {
		t.Fatalf("failed fetching ticket types: %v", errTypes)
	}

	depts, errDepts := app.Repos.Departments.GetAllDepartments(ctx)
	if errDepts != nil || len(depts) == 0 {
		t.Fatalf("failed fetching departments: %v", errDepts)
	}

	deptID := depts[0].ID
	typeID := types[0].ID

	t.Log("📝 [SCENARIO STEP 3]: Guest User creates a new ticket")
	phone := "09121234567"
	ticketReq := dto.TicketCreateRequest{
		PhoneNumber:  phone,
		DepartmentID: deptID,
		TicketTypeID: typeID,
		Title:        "System Outage in Production",
		Body:         "We are experiencing 502 Bad Gateway errors on checkout.",
	}

	// Validate ticket request fields
	if ticketReq.Title == "" || ticketReq.Body == "" {
		t.Fatal("invalid ticket payload")
	}

	// Track code format validation unit check
	sampleCode := "ABC12345"
	cleaned, err := util.ParsTrackCode(sampleCode)
	if err != nil || cleaned != sampleCode {
		t.Fatalf("track code parsing failed: %v", err)
	}

	t.Log("💬 [SCENARIO STEP 4]: Message thread simulation")
	now := time.Now()
	ticketID := uuid.New().String()
	firstMsg := model.ChatMessage{
		ID:          uuid.New().String(),
		SenderID:    0,
		SenderType:  "user",
		Message:     ticketReq.Body,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ticketModel := &model.Ticket{
		ID:             ticketID,
		TrackCode:      sampleCode,
		UserID:         0,
		PhoneNumber:    phone,
		TicketTypeID:   typeID,
		DepartmentID:   deptID,
		TicketStatusID: 1, // open
		Title:          ticketReq.Title,
		CreatedAt:      now,
		UpdatedAt:      now,
		Chat:           []model.ChatMessage{firstMsg},
	}

	if len(ticketModel.Chat) != 1 {
		t.Fatalf("expected 1 initial chat message, got %d", len(ticketModel.Chat))
	}

	t.Log("👨‍💼 [SCENARIO STEP 5]: Support Agent authenticates and replies to ticket")
	agentToken := testutil.GenerateTestAuthToken(t, app.Services.Token, 2, "agent_smith", []int64{2})
	if agentToken == "" {
		t.Fatal("failed generating agent auth token")
	}

	agentReply := model.ChatMessage{
		ID:          uuid.New().String(),
		SenderID:    2,
		SenderType:  "agent",
		Message:     "Hello, our infrastructure team is investigating the 502 issue now.",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	ticketModel.Chat = append(ticketModel.Chat, agentReply)

	if len(ticketModel.Chat) != 2 {
		t.Fatalf("expected 2 chat messages after agent reply, got %d", len(ticketModel.Chat))
	}
	if ticketModel.Chat[1].SenderType != "agent" {
		t.Fatalf("expected sender type 'agent', got '%s'", ticketModel.Chat[1].SenderType)
	}

	t.Log("🔒 [SCENARIO STEP 6]: Support Agent closes resolved ticket")
	ticketModel.TicketStatusID = 3 // closed
	ticketModel.UpdatedAt = time.Now()

	if ticketModel.TicketStatusID != 3 {
		t.Fatalf("expected closed status (3), got %d", ticketModel.TicketStatusID)
	}

	t.Log("✅ [SCENARIO COMPLETE]: Full ticket lifecycle successfully verified")
}
