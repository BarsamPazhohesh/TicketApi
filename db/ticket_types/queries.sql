-- name: GetAllTicketTypes :many
SELECT * FROM ticket_types
WHERE deleted = 0;

-- name: GetAllActiveTicketTypes :many
SELECT * FROM ticket_types
WHERE deleted = 0
AND status != 0;

-- name: AddTicketType :one
INSERT INTO ticket_types (title, description) VALUES (?, ?) RETURNING id;


-- name: CheckTicketTypeByID :one
SELECT COUNT(id) AS exist_of_id
FROM ticket_types
WHERE deleted = 0
AND status != 0
AND id = ?;
