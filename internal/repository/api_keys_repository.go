// Package repository
package repository

import (
	"context"
	"database/sql"
	"ticket-api/internal/db/api_keys"
	"ticket-api/internal/errx"
)

type APIKeysRepository struct {
	queries *api_keys.Queries
}

func NewAPIKeysRepository(queries *api_keys.Queries) *APIKeysRepository {
	return &APIKeysRepository{
		queries: queries,
	}
}

func (repo *APIKeysRepository) AddAPIKey(ctx context.Context, param api_keys.AddApiKeyParams) error {
	err := repo.queries.AddApiKey(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func (repo *APIKeysRepository) GetApiKeyIDByKey(ctx context.Context, key string) (int64, *errx.APIError) {
	id, err := repo.queries.GetActiveAPIKeyID(ctx, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errx.Respond(errx.ErrApiKeyNotFound, err)
		}
		return -1, errx.Respond(errx.ErrInternalServerError, err)
	}

	return id, nil
}
