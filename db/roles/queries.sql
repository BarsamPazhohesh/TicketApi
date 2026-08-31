-- name: IsRoleExist :one
SELECT count(id) as exist_of_id FROM roles
WHERE deleted_at IS NULL
AND deleted = 0
AND status != 0
AND id = ?;

-- name: AddRole :one
INSERT INTO roles(title) VALUES (?) RETURNING id;

-- name: GetAllActiveRoles :many
SELECT id, title, status, created_at, updated_at
FROM roles
WHERE deleted_at IS NULL
AND deleted = 0
AND status != 0;
