// Package repository
package repository

import (
	"context"
	"ticket-api/internal/db/api_handlers"
)

type ApiHandlerRepository struct {
	queries *api_handlers.Queries
}

func NewApiHandlerRepository(queries *api_handlers.Queries) *ApiHandlerRepository {
	return &ApiHandlerRepository{
		queries: queries,
	}
}

func (repo *ApiHandlerRepository) AddApiHandler(ctx context.Context, apiHandler api_handlers.AddApiHandlerParams) (int64, error) {
	apiHandlerId, err := repo.queries.AddApiHandler(ctx, apiHandler)
	if err != nil {
		return -1, err
	}

	return apiHandlerId, nil
}
