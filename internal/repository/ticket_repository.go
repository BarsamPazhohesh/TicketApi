package repository

import (
	"context"
	"strings"
	"ticket-api/internal/config"
	"ticket-api/internal/dto"
	"ticket-api/internal/errx"
	"ticket-api/internal/model"
	"ticket-api/internal/services/storage"
	"ticket-api/internal/util"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TicketRepository handles ticket-related MongoDB operations.
type TicketRepository struct {
	collection *mongo.Collection
	storage    *storage.StorageService
}

// NewTicketRepository initializes a TicketRepository with the "tickets" collection.
// Returns an empty repository if ENABLE_MONGO is 0.
func NewTicketRepository(db *mongo.Database, storage *storage.StorageService) *TicketRepository {
	if !config.Get().Mongo.Enable {
		return &TicketRepository{}
	}
	return &TicketRepository{
		collection: db.Collection(config.Get().Mongo.TicketCollectionName),
		storage:    storage,
	}
}

// CreateTicket inserts a new ticket into MongoDB and returns the ticket ID.
func (r *TicketRepository) CreateTicket(ctx context.Context, ticketDTO *dto.TicketCreateRequest) (*dto.TicketCreateResponse, *errx.APIError) {
	// Parse attachment object names
	attachments, err := util.ParseObjectNames(ticketDTO.Attachments)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}
	ticketDTO.Attachments = attachments

	// Convert DTO to model
	ticket, err := ticketDTO.ToModel(ctx, r.collection)
	if err != nil {
		return nil, errx.Respond(errx.ErrInvalidInput, err)
	}

	// Move temp attachments to ticket folder if first chat has attachments
	if len(ticket.Chat) > 0 && len(ticket.Chat[0].Attachments) > 0 {
		movedAttachments, apiErr := r.storage.MoveTempsFileToTickets(ctx, ticket.ID, attachments)
		if apiErr != nil {
			return nil, apiErr
		}
		ticket.Chat[0].Attachments = movedAttachments
		ticket.AttachmentCount = len(movedAttachments)
	}

	// Insert ticket into MongoDB
	if _, err := r.collection.InsertOne(ctx, ticket); err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	return &dto.TicketCreateResponse{
		ID:        ticket.ID,
		TrackCode: ticket.TrackCode,
	}, nil
}

// GetTicketByID retrieves a single ticket by ID and converts it to TicketRaw.
func (r *TicketRepository) GetTicketByID(ctx context.Context, id string) (*dto.TicketResponse, *errx.APIError) {

	// Validate UUID
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
	// Validate UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		return 0, errx.Respond(errx.ErrBadRequest, err)
	}

	var result struct {
		AttachmentCount int `bson:"attachmentCount"`
	}

	// Use projection to fetch only attachmentCount
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
) (*dto.PagingResponse[dto.TicketResponse], *errx.APIError) {
	cfg := config.Get().TicketConfig

	// Ensure pageSize is within allowed range
	if query.PageSize < cfg.MinPagingSize || query.PageSize > cfg.MaxPagingSize {
		query.PageSize = cfg.DefaultPagingSize
	}
	if query.Page < 1 {
		query.Page = 1
	}

	// Build filter
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

	// Sorting
	allowedSortFields := map[string]bool{
		"createdAt":      true,
		"updatedAt":      true,
		"ticketTypeId":   true,
		"ticketStatusId": true,
		"departmentId":   true,
	}

	sortField := "createdAt"
	if query.OrderBy != "" && allowedSortFields[query.OrderBy] {
		sortField = query.OrderBy
	}

	orderDir := -1
	if strings.ToLower(query.OrderDir) == "asc" {
		orderDir = 1
	}

	skip := (query.Page - 1) * query.PageSize
	// Use bson.D for Sort, bson.M for everything else
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(query.PageSize)).
		SetSort(bson.M{sortField: orderDir}).
		SetProjection(bson.M{"chat": 0}) // exclude chat field

	// Fetch tickets
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}
	defer cursor.Close(ctx)

	var tickets []model.Ticket
	if err = cursor.All(ctx, &tickets); err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	// Total count with cap
	max := cfg.MaxCountingItem
	total, err := r.collection.CountDocuments(ctx, filter, options.Count().SetLimit(max))
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	// Map to DTO
	ticketsDto := make([]dto.TicketResponse, len(tickets))
	for i, ticket := range tickets {
		ticketsDto[i] = *dto.ToTicketResponse(&ticket)
	}

	// Calculate total pages
	totalPages := int(total) / query.PageSize
	if int(total)%query.PageSize != 0 {
		totalPages++
	}

	return &dto.PagingResponse[dto.TicketResponse]{
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
		Total:      total,
		Items:      ticketsDto,
	}, nil
}

func (r *TicketRepository) SetTicketStatus(ctx context.Context, id string, statusId int64) (*dto.TicketResponse, *errx.APIError) {

	// Validate UUID
	uid, err := uuid.Parse(id)

	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	filter := bson.M{"_id": uid.String()}

	update := bson.D{
		{Key: "$set", Value: bson.M{
			"ticketStatusId": statusId,
		}},
		{Key: "$currentDate", Value: bson.M{
			"updatedAt": true,
		}},
	}

	// Options: return the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetProjection(bson.M{"chat": 0})

	var model model.Ticket
	err = r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&model)
	if err != nil {
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	// Convert to DTO
	ticketDTO := dto.ToTicketResponse(&model)

	return ticketDTO, nil
}
