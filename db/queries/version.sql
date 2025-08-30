-- name: GetCurrentVersion :one
SELECT * FROM app_versions
WHERE api_version = ? AND is_current = 1;

-- name: CreateVersion :one
INSERT INTO app_versions (api_version, version, notes, is_current)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: ListVersions :many
SELECT * FROM app_versions
WHERE api_version = ?
ORDER BY release_date DESC;

