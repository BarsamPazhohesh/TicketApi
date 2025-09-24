package repository

import (
	"context"
	"ticket-api/internal/db/roles"
)

type RolesRepository struct {
	queries *roles.Queries
}

func NewRolesRepository(queries *roles.Queries) *RolesRepository {
	return &RolesRepository{
		queries: queries,
	}
}

func (repo *RolesRepository) IsRoleExist(ctx context.Context, roleID int64) (bool, error) {
	_, err := repo.queries.IsRoleExist(ctx, roleID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *RolesRepository) AddRole(ctx context.Context, title string) (int64, error) {
	roleID, err := repo.queries.AddRole(ctx, title)
	if err != nil {
		return -1, err
	}

	return roleID, nil
}
