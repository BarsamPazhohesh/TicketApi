-- name: AddPermission :one
INSERT INTO permissions (name, description)
VALUES (?, ?)
RETURNING id;

-- name: GetAllActivePermissions :many
SELECT id, name, description, status, created_at, updated_at
FROM permissions
WHERE deleted_at IS NULL
AND status != 0;

-- name: GetPermissionByName :one
SELECT id, name, description, status, created_at, updated_at
FROM permissions
WHERE deleted_at IS NULL
AND status != 0
AND name = ?;
