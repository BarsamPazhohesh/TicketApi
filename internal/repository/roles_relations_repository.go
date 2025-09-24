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
	err := repo.roleRelationsQueries.AddApiRoutesToRolesRelation(ctx, param)
	return err
}

func (repo *RolesRelationsRepository) AddTicketTypesToRolesRelation(ctx context.Context, param roles_relations.AddTicketTypesToRolesRelationParams) error {
	err := repo.roleRelationsQueries.AddTicketTypesToRolesRelation(ctx, param)
	return err
}

func (repo *RolesRelationsRepository) AddUsersToRolesRelation(ctx context.Context, param roles_relations.AddUsersToRolesRelationParams) error {
	err := repo.roleRelationsQueries.AddUsersToRolesRelation(ctx, param)
	return err
}
