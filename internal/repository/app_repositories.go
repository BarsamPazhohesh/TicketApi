package repository

import (
	"database/sql"
	"ticket-api/internal/db/api_keys"
	"ticket-api/internal/db/api_routes"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/db/roles"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/db/ticket_priorities"
	"ticket-api/internal/db/ticket_statuses"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/db/users"
	"ticket-api/internal/db/version"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppRepositories struct {
	Ticket           *TicketRepository
	ChatRepository   *ChatRepository
	Version          *VersionRepository
	Roles            *RolesRepository
	Departments      *DepartmentsRepository
	TicketTypes      *TicketTypesRepository
	TicketPriorities *TicketPrioritiesRepository
	APIRoutes        *APIRoutesRepository
	RolesRelations   *RolesRelationsRepository
	Users            *UsersRepository
	TicketStatus     *TicketStatusesRepository
	APIKeys          *APIKeysRepository
}

func NewRepositories(sqldb *sql.DB, mongodb *mongo.Database) *AppRepositories {
	return &AppRepositories{
		Ticket:           NewTicketRepository(mongodb),
		ChatRepository:   NewChatRepository(mongodb),
		Version:          NewVersionRepository(version.New(sqldb)),
		Roles:            NewRolesRepository(roles.New(sqldb)),
		Departments:      NewDepartmentsRepository(departments.New(sqldb)),
		TicketTypes:      NewTicketTypesRepository(ticket_types.New(sqldb)),
		TicketPriorities: NewTicketPrioritiesRepository(ticket_priorities.New(sqldb)),
		APIRoutes:        NewAPIRoutesRepository(api_routes.New(sqldb)),
		APIKeys:          NewAPIKeysRepository(api_keys.New(sqldb)),
		RolesRelations: NewRolesRelationRepository(
			roles_relations.New(sqldb),
			api_keys.New((sqldb)),
			api_routes.New(sqldb)),
		Users:        NewUsersRepository(users.New(sqldb)),
		TicketStatus: NewTicketStatusesRepository(ticket_statuses.New(sqldb)),
	}
}
