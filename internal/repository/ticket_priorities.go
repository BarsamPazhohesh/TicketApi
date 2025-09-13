package repository

import (
	"context"
	"ticket-api/internal/db/ticket_priorities"
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
