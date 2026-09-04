-- name: GetAllSMSMessages :many
SELECT * FROM sms_messages
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: GetSMSMessageByID :one
SELECT * FROM sms_messages
WHERE id = ? AND deleted_at IS NULL;
