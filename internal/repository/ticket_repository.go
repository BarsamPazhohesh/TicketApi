package repository

import (
	"context"
	"errors"
	"ticket-api/internal/dto"
	"ticket-api/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TicketRepository struct {
	collection *mongo.Collection
}

// NewTicketRepository creates a new TicketRepository
func NewTicketRepository(db *mongo.Database) *TicketRepository {
	return &TicketRepository{
		collection: db.Collection("tickets"),
	}
}

// CreateTicket inserts a new ticket and returns the DTO
func (r *TicketRepository) CreateTicket(ctx context.Context, ticketDTO *dto.TicketDTO) (*dto.TicketDTO, error) {
	ticket := dto.FromDTO(ticketDTO) // convert DTO to model

	_, err := r.collection.InsertOne(ctx, ticket)
	if err != nil {
		return nil, err
	}

	return dto.ToDTO(ticket), nil
}

// GetTicket fetches a ticket by ID and returns the DTO
func (r *TicketRepository) GetTicket(ctx context.Context, id string) (*dto.TicketDTO, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ticket ID")
	}

	var ticket model.Ticket
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // not found
		}
		return nil, err
	}

	return dto.ToDTO(&ticket), nil
}

