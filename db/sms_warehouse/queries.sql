-- name: CreateSMSWarehouseRecord :one
INSERT INTO sms_warehouse (
    sms_type_id,
    receiver_phone_number,
    message,
    status
) VALUES (
    ?, ?, ?, ?
) RETURNING *;

-- name: UpdateSMSWarehouseStatus :exec
UPDATE sms_warehouse
SET status = ?,
    updated_at = datetime('now')
WHERE id = ? AND deleted_at IS NULL;

-- name: GetPendingSMSWarehouseRecords :many
SELECT
    sw.id,
    sw.sms_type_id,
    st.title AS sms_type_title,
    sw.receiver_phone_number,
    sw.message,
    sw.status,
    sw.created_at,
    sw.updated_at,
    sw.deleted_at
FROM sms_warehouse sw
INNER JOIN sms_types st ON st.id = sw.sms_type_id AND st.deleted_at IS NULL
WHERE sw.status = 0 AND sw.deleted_at IS NULL
ORDER BY sw.id ASC;

-- name: GetSMSWarehouseByID :one
SELECT
    sw.id,
    sw.sms_type_id,
    st.title AS sms_type_title,
    sw.receiver_phone_number,
    sw.message,
    sw.status,
    sw.created_at,
    sw.updated_at,
    sw.deleted_at
FROM sms_warehouse sw
INNER JOIN sms_types st ON st.id = sw.sms_type_id AND st.deleted_at IS NULL
WHERE sw.id = ? AND sw.deleted_at IS NULL;

