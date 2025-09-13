package repository

import (
	"context"
	"ticket-api/internal/db/roles_relations"
)

type RolesRelationsRepository struct {
	queries *roles_relations.Queries
}

func NewRolesRelationRepository(queries *roles_relations.Queries) *RolesRelationsRepository {
	return &RolesRelationsRepository{
		queries: queries,
	}
}

func (repo *RolesRelationsRepository) AddApiHandlerToRolesRelation(ctx context.Context, param roles_relations.AddApiHandlerToRolesRelationParams) error {
	err := repo.queries.AddApiHandlerToRolesRelation(ctx, param)
	return err
}

func (repo *RolesRelationsRepository) AddTicketTypesToRolesRelation(ctx context.Context, param roles_relations.AddTicketTypesToRolesRelationParams) error {
	err := repo.queries.AddTicketTypesToRolesRelation(ctx, param)
	return err
}

func (repo *RolesRelationsRepository) AddUsersToRolesRelation(ctx context.Context, param roles_relations.AddUsersToRolesRelationParams) error {
	err := repo.queries.AddUsersToRolesRelation(ctx, param)
	return err
}
