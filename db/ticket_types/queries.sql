-- name: GetAllTicketTypes :many
SELECT * FROM ticket_types
WHERE deleted = 0
AND status != 0;

-- name: AddTicketType :one
INSERT INTO ticket_types (title, description) VALUES (?, ?) RETURNING id;