// Package repository
package repository

import (
	"context"
	"ticket-api/internal/db/api_routes"
)

type APIRoutesRepository struct {
	queries *api_routes.Queries
}

func NewAPIRoutesRepository(queries *api_routes.Queries) *APIRoutesRepository {
	return &APIRoutesRepository{
		queries: queries,
	}
}

func (repo *APIRoutesRepository) AddAPIRoute(ctx context.Context, apiRoute api_routes.AddApiRouteParams) (*int64, error) {
	apiRouteID, err := repo.queries.AddApiRoute(ctx, apiRoute)
	if err != nil {
		return nil, err
	}

	return &apiRouteID, nil
}
