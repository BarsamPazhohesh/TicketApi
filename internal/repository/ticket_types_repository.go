package repository

import (
	"context"
	"ticket-api/internal/db/ticket_types"
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
