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

func (repo *RolesRelationsRepository) AddAPIKeysToRolesRelation(ctx context.Context, param roles_relations.AddAPIKeysToRolesRelationParams) error {
	err := repo.roleRelationsQueries.AddAPIKeysToRolesRelation(ctx, param)
	return err
}

func (repo *RolesRelationsRepository) HasRouteAccess(ctx context.Context) (bool, error) {
	// Get ID of apiKey
	apiKeyID, err := repo.apiKeysQueries.GetActiveAPIKeyID(ctx, "SampleKey")
	if err != nil {
		return false, err
	}

	// Get Roles of that ID
	apiKeyToRolesRelationIDs, err := repo.roleRelationsQueries.GetAPIKeyRoleIDs(ctx, apiKeyID)
	if err != nil {
		return false, err
	}

	// Get ID of route
	apiRouteID, err := repo.apiRoutesQueries.GetAPIRouteID(ctx, "/SampleRoute")
	if err != nil {
		return false, err
	}

	// Get Roles of that ID
	apiRouteToRolesRelationIDs, err := repo.roleRelationsQueries.GetAPIRouteRoleIDs(ctx, apiRouteID)
	if err != nil {
		return false, err
	}

	// Compare the two IDs
	routeRolesSet := make(map[int64]struct{}, len(apiRouteToRolesRelationIDs))
	for _, roleID := range apiRouteToRolesRelationIDs {
		routeRolesSet[roleID] = struct{}{}
	}

	// Check if any of the API key's roles exist in the route's roles map.
	for _, apiKeyRoleID := range apiKeyToRolesRelationIDs {
		if _, ok := routeRolesSet[apiKeyRoleID]; ok {
			return true, nil
		}
	}

	return false, nil
}
