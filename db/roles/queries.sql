-- name: CheckRoleExistence :one
SELECT count(id) as exist_of_id FROM roles
WHERE deleted = 0
AND status != 0
AND id = ?;

-- name: InsertRole :one
INSERT INTO roles(title) VALUES (?) RETURNING id;
