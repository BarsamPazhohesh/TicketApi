-- name: AddApiHandlerToRolesRelation :exec
INSERT INTO api_handlers_roles_relation (api_handler_id, role_id) VALUES (?, ?);

-- name: AddUsersToRolesRelation :exec
INSERT INTO users_roles_relation (user_id, role_id) VALUES (?, ?);
