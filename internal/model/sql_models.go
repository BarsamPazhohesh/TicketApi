// Package model
package model

import (
	"ticket-api/internal/db/api_routes"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/db/roles"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/db/sms_messages"
	"ticket-api/internal/db/sms_type_messages_relation"
	"ticket-api/internal/db/sms_types"
	"ticket-api/internal/db/sms_warehouse"
	"ticket-api/internal/db/ticket_statuses"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/db/users"
)

type (
	User                    = users.User
	Role                    = roles.Role
	UsersRolesRelation      = roles_relations.UsersRolesRelation
	Department              = departments.Department
	TicketType              = ticket_types.TicketType
	APIHandler              = api_routes.ApiRoute
	TicketStatus            = ticket_statuses.TicketStatus
	SMSWarehouse            = sms_warehouse.SmsWarehouse
	SMSType                 = sms_types.SmsType
	SMSMessage              = sms_messages.SmsMessage
	SMSTypeMessagesRelation = sms_type_messages_relation.SmsTypeMessagesRelation
)
