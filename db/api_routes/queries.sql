-- name: AddApiRoute :one
INSERT INTO api_routes (route, method, description) VALUES (?, ?, ?) RETURNING id;

-- name: GetAPIRouteID :one
SELECT id FROM api_routes
WHERE deleted = 0
AND status != 0
AND route = ?;
