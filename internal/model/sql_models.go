// Package model
package model

import (
	"ticket-api/internal/db/api_routes"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/db/roles"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/db/ticket_statuses"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/db/users"
	"ticket-api/internal/db/version"
)

type (
	User               = users.User
	Role               = roles.Role
	AppVersion         = version.AppVersion
	UsersRolesRelation = roles_relations.UsersRolesRelation
	Department         = departments.Department
	TicketType         = ticket_types.TicketType
	APIHandler         = api_routes.ApiRoute
	TicketStatus       = ticket_statuses.TicketStatus
)
