package dto

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"ticket-api/internal/model"
	"ticket-api/internal/util"
)

// ---------------------------
// Ticket creation DTOs
// ---------------------------

// TicketCreateRequest represents the payload for creating a new ticket
type TicketCreateRequest struct {
	UserID      int      `json:"userId" binding:"required"`
	Type        int      `json:"type" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Body        string   `json:"body" binding:"required"`
	Attachments []string `json:"attachments,omitempty"`
}

// TicketIDResponse is returned after ticket creation
type TicketIDResponse struct {
	ID string `json:"id"`
}

// ToModel converts a TicketCreateRequest into a model.Ticket
func (dto *TicketCreateRequest) ToModel(ctx context.Context, ticketCollection *mongo.Collection) (*model.Ticket, error) {
	now := time.Now()

	trackCode, err := util.GenerateUniqueTrackCode(ctx, ticketCollection)
	if err != nil {
		return nil, err
	}

	firstMessage := model.ChatMessage{
		ID:          util.GenerateUUID(),
		SenderID:    dto.UserID,
		Message:     dto.Body,
		Attachments: dto.Attachments,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return &model.Ticket{
		ID:             util.GenerateUUID(),
		TrackCode:      trackCode,
		UserID:         dto.UserID,
		TypeID:         dto.Type,
		Title:          dto.Title,
		TicketStatusID: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
		Chat:           []model.ChatMessage{firstMessage},
	}, nil
}

// ---------------------------
// Ticket raw DTO (Mongo document)
// ---------------------------

// TicketRaw represents the ticket data stored in MongoDB
type TicketRaw struct {
	ID             string           `json:"id" bson:"_id"`
	TrackCode      string           `json:"trackCode" bson:"trackCode"`
	UserID         int              `json:"userId" bson:"userId"`
	TypeID         int              `json:"typeId" bson:"typeId"`
	DepartmentID   int              `json:"departmentId" bson:"departmentId"`
	Title          string           `json:"title" bson:"title"`
	TicketStatusID int              `json:"ticketStatusId" bson:"ticketStatusId"`
	CreatedAt      time.Time        `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt" bson:"updatedAt"`
	Chat           []ChatMessageDTO `json:"chat" bson:"chat"`
}

// ToModel converts TicketRaw into model.Ticket
func (r *TicketRaw) ToModel() *model.Ticket {
	chat := make([]model.ChatMessage, len(r.Chat))
	for i, msg := range r.Chat {
		chat[i] = model.ChatMessage{
			ID:          msg.ID,
			SenderID:    msg.SenderID,
			Message:     msg.Message,
			Attachments: msg.Attachments,
			CreatedAt:   msg.CreatedAt,
			UpdatedAt:   msg.UpdatedAt,
		}
	}

	return &model.Ticket{
		ID:             r.ID,
		TrackCode:      r.TrackCode,
		UserID:         r.UserID,
		TypeID:         r.TypeID,
		Title:          r.Title,
		TicketStatusID: r.TicketStatusID,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
		Chat:           chat,
	}
}

// ---------------------------
// TicketResponse DTO (for API)
// ---------------------------

type TicketResponse struct {
	ID             string           `json:"id"`
	TrackCode      string           `json:"trackId"`
	UserID         int              `json:"userId"`
	Username       string           `json:"username"`
	TypeID         int              `json:"typeId"`
	Type           string           `json:"type"`
	DepartmentId   int              `json:"departmentId"`
	DepartmentName string           `json:"departmentName"`
	Priority       int              `json:"priority"`
	Title          string           `json:"title"`
	TicketStatus   string           `json:"ticketStatus"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Chat           []ChatMessageDTO `json:"chat"`
}

func ToTicketRaw(ticket *model.Ticket) *TicketRaw {
	chatDTOs := make([]ChatMessageDTO, len(ticket.Chat))
	for i, msg := range ticket.Chat {
		chatDTOs[i] = ChatMessageDTO{
			ID:          msg.ID,
			SenderID:    msg.SenderID,
			Message:     msg.Message,
			Attachments: msg.Attachments,
			CreatedAt:   msg.CreatedAt,
			UpdatedAt:   msg.UpdatedAt,
		}
	}

	return &TicketRaw{
		ID:             ticket.ID,
		TrackCode:      ticket.TrackCode,
		UserID:         ticket.UserID,
		TypeID:         ticket.TypeID,
		DepartmentID:   ticket.DepartmentID,
		Title:          ticket.Title,
		TicketStatusID: ticket.TicketStatusID,
		CreatedAt:      ticket.CreatedAt,
		UpdatedAt:      ticket.UpdatedAt,
		Chat:           chatDTOs,
	}
}

// TicketToTicketResponse converts a model.Ticket and related data into a TicketResponse
func (r *TicketRaw) ToTicketResponse(
	user *model.User,
	typ *model.TicketType,
	dep *model.Department,
) *TicketResponse {
	chatDTOs := make([]ChatMessageDTO, len(r.Chat))
	for i, msg := range r.Chat {
		chatDTOs[i] = ChatMessageDTO{
			ID:          msg.ID,
			SenderID:    msg.SenderID,
			Message:     msg.Message,
			Attachments: msg.Attachments,
			CreatedAt:   msg.CreatedAt,
			UpdatedAt:   msg.UpdatedAt,
		}
	}

	return &TicketResponse{
		ID:             r.ID,
		UserID:         r.UserID,
		Username:       user.Username,
		TrackCode:      r.TrackCode,
		Priority:       0,
		TypeID:         int(typ.ID),
		Type:           typ.Title,
		DepartmentId:   int(dep.ID),
		DepartmentName: dep.Title,
		Title:          r.Title,
		TicketStatus:   "فعال",
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
		Chat:           chatDTOs,
	}
}
