package repository

import (
	"context"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/model"
	"ticket-api/internal/util"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TicketRepository handles ticket-related MongoDB operations.
type TicketRepository struct {
	collection *mongo.Collection
}

// NewTicketRepository initializes a TicketRepository with the "tickets" collection.
// Returns an empty repository if ENABLE_MONGO is 0.
func NewTicketRepository(db *mongo.Database) *TicketRepository {
	if !config.Get().Mongo.Enable || db == nil {
		return &TicketRepository{}
	}
	return &TicketRepository{
		collection: db.Collection(config.Get().Mongo.TicketCollectionName),
	}
}

func (r *TicketRepository) GetCollection() *mongo.Collection {
	return r.collection
}

// InsertTicket inserts a model.Ticket directly into MongoDB
func (r *TicketRepository) InsertTicket(ctx context.Context, ticket *model.Ticket) *errx.APIError {
	if r.collection == nil {
		return errx.Respond(errx.ErrInternalServerError, nil)
	}
	if _, err := r.collection.InsertOne(ctx, ticket); err != nil {
		return errx.Respond(errx.ErrInternalServerError, err)
	}
	return nil
}

// GetTicketByID retrieves a single ticket by ID and converts it to TicketRaw.
func (r *TicketRepository) GetTicketByID(ctx context.Context, id string) (*dto.TicketResponse, *errx.APIError) {
	if r.collection == nil {
		return nil, errx.Respond(errx.ErrInternalServerError, nil)
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	var ticket model.Ticket
	if err := r.collection.FindOne(ctx, bson.M{"_id": uid.String()}).Decode(&ticket); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errx.Respond(errx.ErrTicketNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	return dto.ToTicketResponse(&ticket), nil
}

func (r *TicketRepository) GetTicketAttachmentCount(ctx context.Context, id string) (int, *errx.APIError) {
	if r.collection == nil {
		return 0, errx.Respond(errx.ErrInternalServerError, nil)
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return 0, errx.Respond(errx.ErrBadRequest, err)
	}

	var result struct {
		AttachmentCount int `bson:"attachmentCount"`
	}

	opts := options.FindOne().SetProjection(bson.M{"attachmentCount": 1})
	if err := r.collection.FindOne(ctx, bson.M{"_id": uid.String()}, opts).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, errx.Respond(errx.ErrTicketNotFound, err)
		}
		return 0, errx.Respond(errx.ErrInternalServerError, err)
	}

	return result.AttachmentCount, nil
}

func (r *TicketRepository) GetTicketByTrackCode(ctx context.Context, trackCode string) (*dto.TicketResponse, *errx.APIError) {
	if r.collection == nil {
		return nil, errx.Respond(errx.ErrInternalServerError, nil)
	}

	code, err := util.ParsTrackCode(trackCode)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	var ticket model.Ticket
	if err := r.collection.FindOne(ctx, bson.M{"trackCode": code}).Decode(&ticket); err != nil {
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

func (r *TicketRepository) GetTickets(
	ctx context.Context,
	query dto.TicketQueryParams,
) (*dto.TicketPagingResponse, *errx.APIError) {
	if r.collection == nil {
		return nil, errx.Respond(errx.ErrInternalServerError, nil)
	}

	cfg := config.Get().TicketConfig

	if query.PageSize < cfg.MinPagingSize || query.PageSize > cfg.MaxPagingSize {
		query.PageSize = cfg.DefaultPagingSize
	}
	if query.Page < 1 {
		query.Page = 1
	}

	filter := bson.M{}
	if query.StatusID != 0 {
		filter["ticketStatusId"] = query.StatusID
	}
	if query.UserID != 0 {
		filter["userId"] = query.UserID
	}
	if query.DepartmentID != 0 {
		filter["departmentId"] = query.DepartmentID
	}
	if query.TicketTypeID != 0 {
		filter["ticketTypeId"] = query.TicketTypeID
	}

	allowedSortFields := map[string]bool{
		"createdAt": true,
		"updatedAt": true,
	}

	sortField := "createdAt"
	if allowedSortFields[query.OrderBy] {
		sortField = query.OrderBy
	}

	sortDir := -1
	if query.OrderDir == "asc" {
		sortDir = 1
	}

	sortOrder := bson.D{{Key: sortField, Value: sortDir}}
	findOptions := options.Find().
		SetSort(sortOrder).
		SetSkip(int64((query.Page - 1) * query.PageSize)).
		SetLimit(int64(query.PageSize))

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
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

	totalPages := int((totalCount + int64(query.PageSize) - 1) / int64(query.PageSize))

	return &dto.TicketPagingResponse{
		Items:      tickets,
		Total:      totalCount,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *TicketRepository) SetTicketStatus(ctx context.Context, ticketID string, statusID int64) (*dto.TicketResponse, *errx.APIError) {
	if r.collection == nil {
		return nil, errx.Respond(errx.ErrInternalServerError, nil)
	}

	uid, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	filter := bson.M{"_id": uid.String()}
	update := bson.M{"$set": bson.M{"ticketStatusId": statusID}}

	res, updateErr := r.collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, updateErr)
	}

	if res.MatchedCount == 0 {
		return nil, errx.Respond(errx.ErrTicketNotFound, nil)
	}

	return r.GetTicketByID(ctx, ticketID)
}
