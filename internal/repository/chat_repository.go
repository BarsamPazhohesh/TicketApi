package repository

import (
	"context"
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChatRepository struct {
	collection *mongo.Collection
}

// NewChatRepository creates a new ChatRepository
func NewChatRepository(db *mongo.Database) *ChatRepository {
	if !config.Get().Mongo.Enable || db == nil {
		return &ChatRepository{}
	}
	return &ChatRepository{
		collection: db.Collection(config.Get().Mongo.TicketCollectionName),
	}
}

// AppendChatMessage adds a chat message model directly to an existing ticket
func (r *ChatRepository) AppendChatMessage(ctx context.Context, ticketID string, message model.ChatMessage) (*dto.ChatMessageDTO, *errx.APIError) {
	if r.collection == nil {
		return nil, errx.Respond(errx.ErrInternalServerError, nil)
	}

	uid, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	update := bson.M{
		"$push": bson.M{"chat": message},
		"$inc":  bson.M{"attachmentCount": len(message.Attachments)},
	}

	res, updateErr := r.collection.UpdateOne(ctx, bson.M{"_id": uid.String()}, update)
	if updateErr != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, updateErr)
	}
	if res.MatchedCount == 0 {
		return nil, errx.Respond(errx.ErrTicketNotFound, errors.New("ticket not found"))
	}

	return &dto.ChatMessageDTO{
		ID:          message.ID,
		SenderID:    message.SenderID,
		Message:     message.Message,
		Attachments: message.Attachments,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   message.UpdatedAt,
	}, nil
}
