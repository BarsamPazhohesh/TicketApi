-- name: AddApiRoutesToRolesRelation :exec
INSERT INTO api_routes_roles_relation (api_route_id, role_id) VALUES (?, ?);

-- name: AddUsersToRolesRelation :exec
INSERT INTO users_roles_relation (user_id, role_id) VALUES (?, ?);

-- name: AddTicketTypesToRolesRelation :exec
INSERT INTO ticket_types_roles_relation (ticket_type_id, role_id) VALUES (?, ?);
