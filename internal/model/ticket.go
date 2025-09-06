package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ticket is the MongoDB model for tickets
type Ticket struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"userId"`      // شناسه کاربر
	Type        string             `bson:"type"`        // نوع تیکت
	Priority    string             `bson:"priority"`    // اولویت
	Title       string             `bson:"title"`       // عنوان
	Body        string             `bson:"body"`        // متن بدنه
	Attachments []string           `bson:"attachments"` // پیوست آرایه آدرس فایل
	Done        bool               `bson:"done"`        // وضعیت انجام شده
	CreatedAt   time.Time          `bson:"createdAt"`   // زمان ایجاد
	UpdatedAt   time.Time          `bson:"updatedAt"`   // زمان بروزرسانی
}
