package repository

import (
	"context"
	"ticket-api/internal/db/ticket_priorities"
	"ticket-api/internal/errx"
)

type TicketPrioritiesRepository struct {
	queries *ticket_priorities.Queries
}

func NewTicketPrioritiesRepository(queries *ticket_priorities.Queries) *TicketPrioritiesRepository {
	return &TicketPrioritiesRepository{
		queries: queries,
	}
}

func (repo *TicketPrioritiesRepository) AddTicketPriority(ctx context.Context, ticketPriority ticket_priorities.AddTicketPriorityParams) error {
	err := repo.queries.AddTicketPriority(ctx, ticketPriority)
	return err
}

func (repo *TicketPrioritiesRepository) GetTicketPriority(ctx context.Context, userID int, ticketTypeID int) (int, *errx.APIError) {
	// TODO: add logic here someone :)
	return -1, nil
}
