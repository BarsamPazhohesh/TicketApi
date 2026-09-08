-- name: GetAllSMSTypes :many
SELECT * FROM sms_types
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: GetSMSTypeByTitle :one
SELECT * FROM sms_types
WHERE title = ? AND deleted_at IS NULL;
