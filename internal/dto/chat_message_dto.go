package dto

import (
	"ticket-api/internal/model"
	"ticket-api/internal/util"
	"time"
)

// ChatMessageDTO represents chat messages in responses
type ChatMessageDTO struct {
	ID          string    `json:"id"`
	SenderID    int       `json:"senderId"`
	Message     string    `json:"message"`
	Attachments []string  `json:"attachments,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ChatMessageResponseID struct {
	ID string `json:"id"`
}

type ChatMessageCreateRequest struct {
	SenderID    int      `json:"senderId"`
	Message     string   `json:"message"`
	Attachments []string `json:"attachments,omitempty"`
}

func (r *ChatMessageCreateRequest) ToModel() *model.ChatMessage {
	now := time.Now()
	return &model.ChatMessage{
		ID:          util.GenerateUUID(),
		SenderID:    r.SenderID,
		Message:     r.Message,
		Attachments: r.Attachments,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
