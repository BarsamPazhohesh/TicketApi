// Package repository
package repository

import (
	"context"
	"ticket-api/internal/db/api_keys"
)

type APIKeysRepository struct {
	queries *api_keys.Queries
}

func NewAPIKeysRepository(queries *api_keys.Queries) *APIKeysRepository {
	return &APIKeysRepository{
		queries: queries,
	}
}

func (repo *APIKeysRepository) AddApiKey(ctx context.Context, param api_keys.AddApiKeyParams) error {
	err := repo.queries.AddApiKey(ctx, param)
	if err != nil {
		return err
	}

	return nil
}
