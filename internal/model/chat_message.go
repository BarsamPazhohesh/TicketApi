package model

import "time"

// ChatMessage represents a single message in a ticket chat
type ChatMessage struct {
	ID          string    `bson:"_id"`
	SenderID    int       `bson:"senderId"`    // شناسه فرستنده
	Message     string    `bson:"message"`     // متن بدنه
	Attachments []string  `bson:"attachments"` // پیوست آرایه آدرس فایل
	CreatedAt   time.Time `bson:"createdAt"`   // زمان ارسال
	UpdatedAt   time.Time `bson:"updatedAt"`   // زمان بروزرسانی
}
