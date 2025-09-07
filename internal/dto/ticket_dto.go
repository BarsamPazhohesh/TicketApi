package dto

import (
	"context"
	"ticket-api/internal/model"
	"ticket-api/internal/util"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// TicketCreateRequest is used when client creates a new ticket
type TicketCreateRequest struct {
	UserID int    `json:"userId" binding:"required"`
	Type   int    `json:"type" binding:"required"`
	Title  string `json:"title" binding:"required"`

	// First message for the ticket (instead of Body/Attachments in Ticket)
	Body        string   `json:"body" binding:"required"`
	Attachments []string `json:"attachments,omitempty"`
}

type TicketIDResponse struct {
	ID string `json:"id"`
}

func (dto *TicketCreateRequest) ToModel(ctx context.Context, ticketCollation *mongo.Collection) (*model.Ticket, error) {
	now := time.Now()
	trackCode, err := util.GenerateUniqueTrackCode(ctx, ticketCollation)
	if err != nil {
		return nil, err
	}

	// create first message
	firstMessage := model.ChatMessage{
		ID:          util.GenerateUUID(),
		SenderID:    dto.UserID,
		Message:     dto.Body,
		Attachments: dto.Attachments,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return &model.Ticket{
		ID:        util.GenerateUUID(),
		TrackCode: trackCode,
		UserID:    dto.UserID,
		Type:      dto.Type,
		Title:     dto.Title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
		Chat:      []model.ChatMessage{firstMessage},
	}, nil
}

// TicketResponse is returned to the client
type TicketResponse struct {
	ID        string           `json:"id"`
	TrackCode string           `json:"trackId"`
	UserID    int              `json:"userId"`
	Type      int              `json:"type"`
	Priority  int              `json:"priority"`
	Title     string           `json:"title"`
	Done      bool             `json:"done"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	Chat      []ChatMessageDTO `json:"chat"`
}

// converts Ticket model to TicketResponse
func TicketToDTO(ticket *model.Ticket) *TicketResponse {
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

	return &TicketResponse{
		ID:        ticket.ID,
		UserID:    ticket.UserID,
		TrackCode: ticket.TrackCode,
		Type:      ticket.Type,
		Title:     ticket.Title,
		Done:      ticket.Done,
		CreatedAt: ticket.CreatedAt,
		UpdatedAt: ticket.UpdatedAt,
		Chat:      chatDTOs,
	}
}
