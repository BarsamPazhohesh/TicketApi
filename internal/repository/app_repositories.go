package repository

import (
	"database/sql"
	"ticket-api/internal/db/api_handlers"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/db/roles"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/db/ticket_priorities"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/db/version"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppRepositories struct {
	Ticket           *TicketRepository
	Version          *VersionRepository
	Roles            *RolesRepository
	Departments      *DepartmentsRepository
	TicketTypes      *TicketTypesRepository
	TicketPriorities *TicketPrioritiesRepository
	ApiHandlers      *ApiHandlerRepository
	RolesRelations   *RolesRelationsRepository
}

func NewRepositories(sqldb *sql.DB, mongodb *mongo.Database) *AppRepositories {
	return &AppRepositories{
		Ticket:           NewTicketRepository(mongodb),
		Version:          NewVersionRepository(version.New(sqldb)),
		Roles:            NewRolesRepository(roles.New(sqldb)),
		Departments:      NewDepartmentsRepository(departments.New(sqldb)),
		TicketTypes:      NewTicketTypesRepository(ticket_types.New(sqldb)),
		TicketPriorities: NewTicketPrioritiesRepository(ticket_priorities.New(sqldb)),
		ApiHandlers:      NewApiHandlerRepository(api_handlers.New(sqldb)),
		RolesRelations:   NewRolesRelationRepository(roles_relations.New(sqldb)),
	}
}
