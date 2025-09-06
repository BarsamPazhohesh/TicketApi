-- name: AddApiHandlerToRolesRelation :exec
INSERT INTO api_handlers_roles_relation (api_handler_id, role_id) VALUES (?, ?);
