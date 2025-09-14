package repository

import (
	"context"
	"ticket-api/internal/dto"
	"ticket-api/internal/env"
	"ticket-api/internal/errx"
	"ticket-api/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
func (r *TicketRepository) CreateTicket(ctx context.Context, ticketDTO *dto.TicketCreateRequest) (*dto.IDResponse[string], *errx.APIError) {
	ticket, err := ticketDTO.ToModel(ctx, r.collection)
	if err != nil {
		return nil, errx.Respond(errx.ErrInvalidInput, err)
	}

	if _, err := r.collection.InsertOne(ctx, ticket); err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	return &dto.IDResponse[string]{ID: ticket.ID}, nil
}

// GetTicket retrieves a single ticket by ID and converts it to TicketRaw.
func (r *TicketRepository) GetTicket(ctx context.Context, id string) (*dto.TicketResponse, *errx.APIError) {
	var ticket model.Ticket
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&ticket); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errx.Respond(errx.ErrTicketNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	return dto.ToTicketResponse(&ticket), nil
}

// GetAllTickets retrieves all tickets for a specific user and converts them to TicketRaw.
func (r *TicketRepository) GetAllTickets(ctx context.Context, userID int) ([]dto.TicketResponse, *errx.APIError) {
	if r.collection == nil {
		return nil, errx.Respond(errx.ErrInternalServerError, nil)
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	defer cursor.Close(ctx)

	var tickets []dto.TicketResponse
	for cursor.Next(ctx) {
		var t model.Ticket
		if err := cursor.Decode(&t); err != nil {
			return nil, errx.Respond(errx.ErrInternalServerError, err)
		}
		tickets = append(tickets, *dto.ToTicketResponse(&t))
	}

	return tickets, nil
}
