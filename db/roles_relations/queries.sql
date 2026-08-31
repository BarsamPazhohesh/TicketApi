-- name: AddApiRoutesToRolesRelation :exec
INSERT INTO api_routes_roles_relation (api_route_id, role_id) VALUES (?, ?);

-- name: AddUsersToRolesRelation :exec
INSERT INTO users_roles_relation (user_id, role_id) VALUES (?, ?);

-- name: AddTicketTypesToRolesRelation :exec
INSERT INTO ticket_types_roles_relation (ticket_type_id, role_id) VALUES (?, ?);

-- name: AddAPIKeysToRolesRelation :exec
INSERT INTO api_keys_roles_relation (api_key_id, role_id) VALUES (?, ?);

-- name: AddRolePermissionRelation :exec
INSERT INTO roles_permissions_relation (role_id, permission_id) VALUES (?, ?);

-- name: AddApiRoutePermissionRelation :exec
INSERT INTO api_routes_permissions_relation (api_route_id, permission_id) VALUES (?, ?);

-- name: GetAPIKeyRoleIDs :many
SELECT role_id FROM api_keys_roles_relation
WHERE deleted_at IS NULL
AND deleted = 0
AND status != 0
AND api_key_id = ?;

-- name: GetAPIRouteRoleIDs :many
SELECT role_id FROM api_routes_roles_relation
WHERE deleted_at IS NULL
AND deleted = 0
AND status != 0
AND api_route_id = ?;

-- name: GetUserRoleIDs :many
SELECT role_id FROM users_roles_relation
WHERE deleted_at IS NULL
AND deleted = 0
AND status != 0
AND user_id = ?;

-- name: GetAllRolesWithPermissions :many
SELECT rpr.role_id, p.name AS permission_name
FROM roles_permissions_relation rpr
JOIN permissions p ON p.id = rpr.permission_id
WHERE rpr.deleted_at IS NULL
AND rpr.status != 0
AND p.deleted_at IS NULL
AND p.status != 0;

-- name: GetAllRoutesWithPermissions :many
SELECT ar.route, ar.method, ar.status, p.name AS permission_name
FROM api_routes ar
LEFT JOIN api_routes_permissions_relation arpr ON arpr.api_route_id = ar.id AND arpr.deleted_at IS NULL AND arpr.status != 0
LEFT JOIN permissions p ON p.id = arpr.permission_id AND p.deleted_at IS NULL AND p.status != 0
WHERE (ar.deleted_at IS NULL)
AND ar.deleted = 0;
