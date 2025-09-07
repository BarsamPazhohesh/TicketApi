package dto

import "time"

// ChatMessageDTO represents chat messages in responses
type ChatMessageDTO struct {
	ID          string    `json:"id"`
	SenderID    int       `json:"senderId"`
	Message     string    `json:"message"`
	Attachments []string  `json:"attachments,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
