package dto

import (
	"ticket-api/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TicketDTO is used to transfer ticket data in requests/responses
type TicketDTO struct {
	UserID      string    `json:"userId"`      // شناسه کاربر
	Type        string    `json:"type"`        // نوع تیکت
	Priority    string    `json:"priority"`    // اولویت
	Title       string    `json:"title"`       // عنوان
	Body        string    `json:"body"`        // متن بدنه
	Attachments []string  `json:"attachments"` // پیوست آرایه آدرس فایل
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// FromDTO converts TicketDTO to Ticket model
func FromDTO(dto *TicketDTO) *model.Ticket {
	now := time.Now()
	return &model.Ticket{
		ID:          primitive.NewObjectID(),
		UserID:      dto.UserID,
		Type:        dto.Type,
		Priority:    dto.Priority,
		Title:       dto.Title,
		Body:        dto.Body,
		Attachments: dto.Attachments,
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ToDTO converts Ticket model to TicketDTO
func ToDTO(ticket *model.Ticket) *TicketDTO {
	return &TicketDTO{
		UserID:      ticket.UserID,
		Type:        ticket.Type,
		Priority:    ticket.Priority,
		Title:       ticket.Title,
		Body:        ticket.Body,
		Attachments: ticket.Attachments,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}
}
