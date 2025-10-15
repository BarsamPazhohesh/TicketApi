package repository

import (
	"context"
	"errors"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/services/storage"
	"ticket-api/internal/util"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChatRepository struct {
	collection *mongo.Collection
	storage    *storage.StorageService
}

// NewChatRepository creates a new ChatRepository
func NewChatRepository(db *mongo.Database, storage *storage.StorageService) *ChatRepository {
	if !config.Get().Mongo.Enable {
		return &ChatRepository{}
	}
	return &ChatRepository{
		collection: db.Collection(config.Get().Mongo.TicketCollectionName),
		storage:    storage,
	}
}

// CreateChatMessageForTicket adds a chat message to an existing ticket
func (r *ChatRepository) CreateChatMessageForTicket(ctx context.Context, ticketID string, message *dto.ChatMessageCreateRequest) (*dto.ChatMessageDTO, *errx.APIError) {

	// Validate UUID
	uid, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	model := message.ToModel()
	attachments, err := util.ParseObjectNames(model.Attachments)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	attachments, err = r.storage.MoveTempsFileToTickets(ctx, ticketID, attachments)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	model.Attachments = attachments
	update := bson.M{
		"$push": bson.M{"chat": model},
		"$inc":  bson.M{"attachmentCount": len(attachments)},
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": uid}, update)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	if res.MatchedCount == 0 {
		return nil, errx.Respond(errx.ErrTicketNotFound, errors.New("ticket not found"))
	}
	return &dto.ChatMessageDTO{
		ID:          model.ID,
		SenderID:    model.SenderID,
		Message:     model.Message,
		Attachments: model.Attachments,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}
