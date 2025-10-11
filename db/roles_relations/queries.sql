-- name: AddApiRoutesToRolesRelation :exec
INSERT INTO api_routes_roles_relation (api_route_id, role_id) VALUES (?, ?);

-- name: AddUsersToRolesRelation :exec
INSERT INTO users_roles_relation (user_id, role_id) VALUES (?, ?);

-- name: AddTicketTypesToRolesRelation :exec
INSERT INTO ticket_types_roles_relation (ticket_type_id, role_id) VALUES (?, ?);

-- name: AddAPIKeysToRolesRelation :exec
INSERT INTO api_keys_roles_relation (api_key_id, role_id) VALUES (?, ?);

-- name: GetAPIKeyRoleIDs :many
SELECT role_id FROM api_keys_roles_relation
WHERE deleted = 0
AND status != 0
AND api_key_id = ?;

-- name: GetAPIRouteRoleIDs :many
SELECT role_id FROM api_routes_roles_relation
WHERE deleted = 0
AND status != 0
AND api_route_id = ?;
