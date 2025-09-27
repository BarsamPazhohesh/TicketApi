package dto

import (
	"context"
	"ticket-api/internal/model"
	"ticket-api/internal/util"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ---------------------------
// Ticket creation DTOs
// ---------------------------

// TicketCreateRequest represents the payload for creating a new ticket
type TicketCreateRequest struct {
	UserID         int64    `json:"userId" binding:"required"`
	TicketTypeID   int64    `json:"ticketTypeID" binding:"required"`
	TicketStatusID int64    `json:"ticketStatusID" binding:"required"`
	DepartmentID   int64    `json:"departmentId" binding:"required"`
	Title          string   `json:"title" binding:"required"`
	Body           string   `json:"body" binding:"required"`
	Attachments    []string `json:"attachments,omitempty"`
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
		TicketTypeID:   dto.TicketTypeID,
		DepartmentID:   dto.DepartmentID,
		TicketStatusID: dto.TicketStatusID,
		Title:          dto.Title,
		CreatedAt:      now,
		UpdatedAt:      now,
		Chat:           []model.ChatMessage{firstMessage},
	}, nil
}

// ---------------------------
// Ticket Response DTO (Mongo document)
// ---------------------------

// TicketResponse represents the ticket data stored in MongoDB
type TicketResponse struct {
	ID             string           `json:"id" bson:"_id"`
	TrackCode      string           `json:"trackCode" bson:"trackCode"`
	UserID         int64            `json:"userId" bson:"userId"`
	TicketTypeID   int64            `json:"ticketTypeId" bson:"typeId"`
	DepartmentID   int64            `json:"departmentId" bson:"departmentId"`
	Title          string           `json:"title" bson:"title"`
	TicketStatusID int64            `json:"ticketStatusId" bson:"ticketStatusId"`
	CreatedAt      time.Time        `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt" bson:"updatedAt"`
	Chat           []ChatMessageDTO `json:"chat" bson:"chat"`
}

// ToModel converts TicketRaw into model.Ticket
func (r *TicketResponse) ToModel() *model.Ticket {
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
		TicketTypeID:   r.TicketTypeID,
		TicketStatusID: r.TicketStatusID,
		DepartmentID:   r.DepartmentID,
		Title:          r.Title,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
		Chat:           chat,
	}
}

// TicketFullResponse DTO (for API)
type TicketFullResponse struct {
	ID             string           `json:"id"`
	TrackCode      string           `json:"trackId"`
	UserID         int64            `json:"userId"`
	Username       string           `json:"username"`
	TicketTypeID   int64            `json:"ticketTypeId"`
	TicketType     string           `json:"ticketType"`
	DepartmentID   int64            `json:"departmentId"`
	DepartmentName string           `json:"departmentName"`
	Priority       int              `json:"priority"`
	Title          string           `json:"title"`
	TicketStatus   string           `json:"ticketStatus"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Chat           []ChatMessageDTO `json:"chat"`
}

func ToTicketResponse(ticket *model.Ticket) *TicketResponse {
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
		ID:             ticket.ID,
		TrackCode:      ticket.TrackCode,
		UserID:         ticket.UserID,
		TicketTypeID:   ticket.TicketTypeID,
		DepartmentID:   ticket.DepartmentID,
		Title:          ticket.Title,
		TicketStatusID: ticket.TicketStatusID,
		CreatedAt:      ticket.CreatedAt,
		UpdatedAt:      ticket.UpdatedAt,
		Chat:           chatDTOs,
	}
}

type TicketByIDRequestDTO struct {
	ID string `json:"id" binding:"required,uuid"` // assuming UUID
}

type TicketCreateResponse struct {
	ID        string `json:"id"`
	TrackCode string `json:"trackCode"`
}

type TicketByTrackCodeRequestDTO struct {
	TrackCode string `json:"trackCode" binding:"required"`
	Username  string `json:"username" binding:"required"`
}

type TicketQueryParams struct {
	Page     int `json:"page,omitempty"`      // page number
	PageSize int `json:"page_size,omitempty"` // items per page

	StatusID     int64 `json:"status,omitempty"`  // optional filter
	UserID       int64 `json:"user_id,omitempty"` // optional filter
	DepartmentID int64 `json:"departmentId,omitempty"`

	OrderBy  string `json:"order_by,omitempty"`  // field to order by
	OrderDir string `json:"order_dir,omitempty"` // asc or desc
}
