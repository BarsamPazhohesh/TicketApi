package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	apperror "ticket-api/internal/apperror"
	"ticket-api/internal/dto"
	"ticket-api/internal/env"
	"ticket-api/internal/model"
)

// TicketRepository handles ticket-related MongoDB operations.
type TicketRepository struct {
	collection *mongo.Collection
}

// NewTicketRepository initializes a TicketRepository with the "tickets" collection.
// Returns an empty repository if ENABLE_MONGO is 0.
func NewTicketRepository(db *mongo.Database) *TicketRepository {
	if env.GetEnvInt("ENABLE_MONGO", 0) == 0 {
		return &TicketRepository{}
	}
	return &TicketRepository{
		collection: db.Collection("tickets"),
	}
}

// CreateTicket inserts a new ticket into MongoDB and returns the ticket ID.
func (r *TicketRepository) CreateTicket(ctx context.Context, ticketDTO *dto.TicketCreateRequest) (*dto.TicketIDResponse, *apperror.APIError) {
	ticket, err := ticketDTO.ToModel(ctx, r.collection)
	if err != nil {
		return nil, apperror.Respond(apperror.ErrInvalidInput, err)
	}

	if _, err := r.collection.InsertOne(ctx, ticket); err != nil {
		return nil, apperror.Respond(apperror.ErrInternalServerError, err)
	}

	return &dto.TicketIDResponse{ID: ticket.ID}, nil
}

// GetTicket retrieves a single ticket by ID and converts it to TicketRaw.
func (r *TicketRepository) GetTicket(ctx context.Context, id string) (*dto.TicketRaw, *apperror.APIError) {
	var ticket model.Ticket
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&ticket); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.Respond(apperror.ErrTicketNotFound, err)
		}
		return nil, apperror.Respond(apperror.ErrInternalServerError, err)
	}

	return dto.ToTicketRaw(&ticket), nil
}

// GetAllTickets retrieves all tickets for a specific user and converts them to TicketRaw.
func (r *TicketRepository) GetAllTickets(ctx context.Context, userID int) ([]dto.TicketRaw, *apperror.APIError) {
	if r.collection == nil {
		return nil, apperror.Respond(apperror.ErrInternalServerError, nil)
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, apperror.Respond(apperror.ErrInternalServerError, err)
	}
	defer cursor.Close(ctx)

	var tickets []dto.TicketRaw
	for cursor.Next(ctx) {
		var t model.Ticket
		if err := cursor.Decode(&t); err != nil {
			return nil, apperror.Respond(apperror.ErrInternalServerError, err)
		}
		tickets = append(tickets, *dto.ToTicketRaw(&t))
	}

	return tickets, nil
}
