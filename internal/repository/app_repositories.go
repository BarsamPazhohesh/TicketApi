package repository

import (
	"database/sql"
	"ticket-api/internal/db/api_keys"
	"ticket-api/internal/db/api_routes"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/db/roles"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/db/sms_type_messages_relation"
	"ticket-api/internal/db/sms_warehouse"
	"ticket-api/internal/db/ticket_priorities"
	"ticket-api/internal/db/ticket_statuses"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/db/users"
	"ticket-api/internal/services/cache"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppRepositories struct {
	Ticket           *TicketRepository
	ChatRepository   *ChatRepository
	Roles            *RolesRepository
	Departments      *DepartmentsRepository
	TicketTypes      *TicketTypesRepository
	TicketPriorities *TicketPrioritiesRepository
	APIRoutes        *APIRoutesRepository
	RolesRelations   *RolesRelationsRepository
	Users            *UsersRepository
	TicketStatus     *TicketStatusesRepository
	APIKeys          *APIKeysRepository
	SMSWarehouse     *SMSWarehouseRepository
	SMSTypeMessages  *sms_type_messages_relation.Queries
}

func NewRepositories(sqldb *sql.DB, mongodb *mongo.Database, redis *redis.Client) *AppRepositories {
	cacheSvc := cache.NewCacheService(redis)

	return &AppRepositories{
		Ticket:           NewTicketRepository(mongodb),
		ChatRepository:   NewChatRepository(mongodb),
		Roles:            NewRolesRepository(roles.New(sqldb)),
		Departments:      NewDepartmentsRepository(departments.New(sqldb), cacheSvc),
		TicketTypes:      NewTicketTypesRepository(ticket_types.New(sqldb), cacheSvc),
		TicketPriorities: NewTicketPrioritiesRepository(ticket_priorities.New(sqldb)),
		APIRoutes:        NewAPIRoutesRepository(api_routes.New(sqldb)),
		APIKeys:          NewAPIKeysRepository(api_keys.New(sqldb)),
		RolesRelations: NewRolesRelationRepository(
			roles_relations.New(sqldb),
			api_keys.New(sqldb),
			api_routes.New(sqldb)),
		Users:           NewUsersRepository(users.New(sqldb)),
		TicketStatus:    NewTicketStatusesRepository(ticket_statuses.New(sqldb), cacheSvc),
		SMSWarehouse:    NewSMSWarehouseRepository(sms_warehouse.New(sqldb)),
		SMSTypeMessages: sms_type_messages_relation.New(sqldb),
	}
}
