package model

import (
	"ticket-api/internal/db/api_handlers"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/db/roles"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/db/users"
	"ticket-api/internal/db/version"
)

type User = users.User
type Role = roles.Role
type AppVersion = version.AppVersion
type UsersRolesRelation = roles_relations.UsersRolesRelation
type Department = departments.Department
type TicketType = ticket_types.TicketType
type ApiHandler = api_handlers.ApiHandler
