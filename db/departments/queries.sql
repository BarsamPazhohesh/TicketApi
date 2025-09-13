-- name: GetAllDepartments :many
SELECT * FROM departments
WHERE deleted = 0;

-- name: GetAllActiveDepartments :many
SELECT * FROM departments
WHERE deleted = 0
AND status != 0;

-- name: AddDepartment :one
INSERT INTO departments (title, description) VALUES (?, ?) RETURNING id;


-- name: GetDepartmentByID :one
SELECT * FROM departments WHERE deleted = ? AND status = ?;
