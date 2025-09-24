package repository

import (
	"context"
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChatRepository struct {
	collection *mongo.Collection
}

// NewChatRepository creates a new ChatRepository
func NewChatRepository(db *mongo.Database) *ChatRepository {
	if !config.Get().Mongo.Enable {
		return &ChatRepository{}
	}
	return &ChatRepository{
		collection: db.Collection(config.Get().Mongo.TicketCollectionName),
	}
}

// CreateChatMessageForTicket adds a chat message to an existing ticket
func (r *TicketRepository) CreateChatMessageForTicket(ctx context.Context, ticketID string, message *dto.ChatMessageCreateRequest) (*dto.ChatMessageResponseID, *errx.APIError) {
	model := message.ToModel()
	update := bson.M{
		"$push": bson.M{"chat": model},
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": ticketID}, update)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	if res.MatchedCount == 0 {
		return nil, errx.Respond(errx.ErrTicketNotFound, errors.New("ticket not found"))
	}
	return &dto.ChatMessageResponseID{ID: model.ID}, nil
}
