package repository

import (
	"context"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/errx"
)

type TicketTypesRepository struct {
	queries *ticket_types.Queries
}

func NewTicketTypesRepository(queries *ticket_types.Queries) *TicketTypesRepository {
	return &TicketTypesRepository{
		queries: queries,
	}
}

func (repo *TicketTypesRepository) AddTicketType(ctx context.Context, ticketType ticket_types.AddTicketTypeParams) (int64, error) {
	ticketTypeId, err := repo.queries.AddTicketType(ctx, ticketType)
	if err != nil {
		return -1, err
	}

	return ticketTypeId, nil
}

func (repo *TicketTypesRepository) GetAllTicketTypes(ctx context.Context) ([]ticket_types.TicketType, error) {
	return repo.queries.GetAllTicketTypes(ctx)
}

func (repo *TicketTypesRepository) GetAllActiveTicketTypes(ctx context.Context) ([]ticket_types.TicketType, error) {
	return repo.queries.GetAllActiveTicketTypes(ctx)
}

func (repo *TicketTypesRepository) IsTicketTypeExits(ctx context.Context, typeID int64) (bool, *errx.APIError) {
	count, err := repo.queries.CheckTicketTypeByID(ctx, typeID)
	if err != nil {
		return false, errx.Respond(errx.ErrInternalServerError, err)
	}
	return count != 0, nil
}
