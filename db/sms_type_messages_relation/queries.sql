-- name: GetAllSMSTypeMessageRelations :many
SELECT
    st.id AS sms_type_id,
    st.title AS sms_type_title,
    sm.id AS sms_message_id,
    sm.title AS sms_message_title,
    sm.body AS sms_message_body
FROM sms_type_messages_relation stmr
INNER JOIN sms_types st ON st.id = stmr.sms_type_id AND st.deleted_at IS NULL
INNER JOIN sms_messages sm ON sm.id = stmr.sms_message_id AND sm.deleted_at IS NULL
WHERE stmr.deleted_at IS NULL;

-- name: GetMessageBySMSTypeTitle :one
SELECT
    sm.body
FROM sms_type_messages_relation stmr
INNER JOIN sms_types st ON st.id = stmr.sms_type_id AND st.deleted_at IS NULL
INNER JOIN sms_messages sm ON sm.id = stmr.sms_message_id AND sm.deleted_at IS NULL
WHERE st.title = ? AND stmr.deleted_at IS NULL
LIMIT 1;
