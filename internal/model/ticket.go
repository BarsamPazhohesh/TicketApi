package model

import (
	"time"
)

// Ticket is the MongoDB model for tickets
type Ticket struct {
	ID             string        `bson:"_id"`            // Unique ticket ID (UUID)
	UserID         int           `bson:"userId"`         // ID of the user who created the ticket
	DepartmentID   int           `bson:"departmentId"`   // Department of the user
	TicketTypeID   int           `bson:"ticketTypeId"`   // Type/category of the ticket
	TicketStatusID int           `bson:"ticketStatusId"` // Current status (open, closed, etc.)
	Title          string        `bson:"title"`          // Short descriptive title
	TrackCode      string        `bson:"trackCode"`      // 8-char code shown to user
	CreatedAt      time.Time     `bson:"createdAt"`      // Ticket creation timestamp
	UpdatedAt      time.Time     `bson:"updatedAt"`      // Last update timestamp
	Chat           []ChatMessage `bson:"chat"`           // Conversation messages for this ticket
}
