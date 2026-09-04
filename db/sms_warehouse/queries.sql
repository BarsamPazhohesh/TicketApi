-- name: CreateSMSWarehouseRecord :one
INSERT INTO sms_warehouse (
    receiver_phone_number,
    message,
    status
) VALUES (
    ?, ?, ?
) RETURNING *;

-- name: UpdateSMSWarehouseStatus :exec
UPDATE sms_warehouse
SET status = ?,
    updated_at = datetime('now')
WHERE id = ? AND deleted_at IS NULL;

-- name: GetPendingSMSWarehouseRecords :many
SELECT * FROM sms_warehouse
WHERE status = 0 AND deleted_at IS NULL
ORDER BY id ASC;

-- name: GetSMSWarehouseByID :one
SELECT * FROM sms_warehouse
WHERE id = ? AND deleted_at IS NULL;
