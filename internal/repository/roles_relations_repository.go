package repository

import (
	"context"
	"ticket-api/internal/db/api_keys"
	"ticket-api/internal/db/api_routes"
	"ticket-api/internal/db/roles_relations"
)

type RolesRelationsRepository struct {
	roleRelationsQueries *roles_relations.Queries
	apiKeysQueries       *api_keys.Queries
	apiRoutesQueries     *api_routes.Queries
}

func NewRolesRelationRepository(
	roleRelationsQueries *roles_relations.Queries,
	apiKeysQueries *api_keys.Queries,
	apiRoutesQueries *api_routes.Queries,
) *RolesRelationsRepository {
	return &RolesRelationsRepository{
		roleRelationsQueries: roleRelationsQueries,
		apiKeysQueries:       apiKeysQueries,
		apiRoutesQueries:     apiRoutesQueries,
	}
}

func (repo *RolesRelationsRepository) AddAPIHandlerToRolesRelation(ctx context.Context, param roles_relations.AddApiRoutesToRolesRelationParams) error {
	return repo.roleRelationsQueries.AddApiRoutesToRolesRelation(ctx, param)
}

func (repo *RolesRelationsRepository) AddTicketTypesToRolesRelation(ctx context.Context, param roles_relations.AddTicketTypesToRolesRelationParams) error {
	return repo.roleRelationsQueries.AddTicketTypesToRolesRelation(ctx, param)
}

func (repo *RolesRelationsRepository) AddUsersToRolesRelation(ctx context.Context, param roles_relations.AddUsersToRolesRelationParams) error {
	return repo.roleRelationsQueries.AddUsersToRolesRelation(ctx, param)
}

func (repo *RolesRelationsRepository) AddAPIKeysToRolesRelation(ctx context.Context, param roles_relations.AddAPIKeysToRolesRelationParams) error {
	return repo.roleRelationsQueries.AddAPIKeysToRolesRelation(ctx, param)
}

func (repo *RolesRelationsRepository) GetAllRolesWithPermissions(ctx context.Context) ([]roles_relations.GetAllRolesWithPermissionsRow, error) {
	return repo.roleRelationsQueries.GetAllRolesWithPermissions(ctx)
}

func (repo *RolesRelationsRepository) GetAllRoutesWithPermissions(ctx context.Context) ([]roles_relations.GetAllRoutesWithPermissionsRow, error) {
	return repo.roleRelationsQueries.GetAllRoutesWithPermissions(ctx)
}

func (repo *RolesRelationsRepository) GetUserRoleIDs(ctx context.Context, userID int64) ([]int64, error) {
	return repo.roleRelationsQueries.GetUserRoleIDs(ctx, userID)
}
