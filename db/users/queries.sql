-- name: GetUserByUsername :one
SELECT id FROM users
WHERE deleted = 0
AND status != 0
AND username = ?;

-- name: CheckUserByID :one
SELECT count(id) exist_of_id FROM users
WHERE deleted = 0
AND status != 0
AND id = ?;

-- name: CreateUser :one
INSERT INTO users (username, department_id) VALUES (?, ?) RETURNING id;

-- name: GetUser :one
SELECT u.id, u.username, u.department_id, u.created_at, u.updated_at, u.status, u.deleted
FROM users u
WHERE u.id = ? AND u.deleted = 0;
