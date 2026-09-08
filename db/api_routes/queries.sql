-- name: AddApiRoute :one
INSERT INTO api_routes (route, method, description) VALUES (?, ?, ?) RETURNING id;

-- name: GetAPIRouteID :one
SELECT id FROM api_routes
WHERE deleted_at IS NULL
AND deleted = 0
AND status != 0
AND route = ?;

-- name: GetAllActiveApiRoutes :many
SELECT id, route, method, description, status, created_at, updated_at
FROM api_routes
WHERE deleted_at IS NULL
AND deleted = 0
AND status != 0;
