package repository

import (
	"context"
	"database/sql"
	"errors"
	"ticket-api/internal/db/ticket_statuses"
	"ticket-api/internal/errx"
)

type TicketStatusesRepository struct {
	queries *ticket_statuses.Queries
}

func NewTicketStatusesRepository(queries *ticket_statuses.Queries) *TicketStatusesRepository {
	return &TicketStatusesRepository{
		queries: queries,
	}
}

func (repo *TicketStatusesRepository) AddTicketStatus(ctx context.Context, param ticket_statuses.AddTicketStatusParams) *errx.APIError {
	err := repo.queries.AddTicketStatus(ctx, param)
	if err != nil {
		return errx.Respond(errx.ErrInternalServerError, err)
	}
	return nil
}

func (repo *TicketStatusesRepository) GetAllActiveTicketStatuses(ctx context.Context) ([]ticket_statuses.TicketStatus, *errx.APIError) {
	data, err := repo.queries.GetAllActiveTicketStatuses(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.Respond(errx.ErrTicketStatusNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	return data, nil
}

func (repo *TicketStatusesRepository) GetActiveTicketStatusByID(ctx context.Context, ID int64) (*ticket_statuses.TicketStatus, *errx.APIError) {
	data, err := repo.queries.GetActiveTicketStatusById(ctx, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.Respond(errx.ErrTicketStatusNotFound, err)
		}
		return nil, errx.Respond(errx.ErrInternalServerError, err)
	}

	return &data, nil
}
