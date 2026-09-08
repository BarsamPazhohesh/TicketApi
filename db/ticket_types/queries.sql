-- name: GetAllTicketTypes :many
SELECT * FROM ticket_types
WHERE deleted_at IS NULL;

-- name: GetAllActiveTicketTypes :many
SELECT * FROM ticket_types
WHERE deleted_at IS NULL
AND status != 0;

-- name: AddTicketType :one
INSERT INTO ticket_types (title, description) VALUES (?, ?) RETURNING id;

-- name: CheckTicketTypeByID :one
SELECT COUNT(id) AS exist_of_id
FROM ticket_types
WHERE deleted_at IS NULL
AND status != 0
AND id = ?;
