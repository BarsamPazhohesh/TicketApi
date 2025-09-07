package model

import (
	"time"
)

// Ticket is the MongoDB model for tickets
type Ticket struct {
	ID        string        `bson:"_id"`       // UUID generated in Go
	TrackCode string        `bson:"TrackCode"` // 8-char user-facing code
	UserID    int           `bson:"userId"`
	Type      int           `bson:"type"`
	Title     string        `bson:"title"`
	Done      bool          `bson:"done"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
	Chat      []ChatMessage `bson:"chat"`
}
