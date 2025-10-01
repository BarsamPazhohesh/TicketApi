-- name: GetUserByUsername :one
SELECT * FROM users
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


-- name: CreateUserWithPassword :one
INSERT INTO users (username, password, department_id) VALUES (?, ?, ?) RETURNING id;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? 
AND deleted = 0 
AND status != 0;

-- name: GetUsersByIDs :many
SELECT * FROM users
WHERE id IN (sqlc.slice('ids'))
AND deleted = 0 
AND status != 0;




