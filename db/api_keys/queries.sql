-- name: AddApiKey :exec
INSERT INTO api_keys (key, description) VALUES ( ?, ? );

-- name: GetActiveAPIKeyID :one
SELECT id FROM api_keys
WHERE deleted_at IS NULL
AND status != 0
AND key = ?;
