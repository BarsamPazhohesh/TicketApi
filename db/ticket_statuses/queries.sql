-- name: AddTicketStatus :exec
INSERT INTO ticket_statuses (title, description) VALUES (?, ?);

-- name: GetAllActiveTicketStatuses :many
SELECT * FROM ticket_statuses
WHERE deleted = 0
AND status != 0;

-- name: GetTicketStatusById :one
SELECT * FROM ticket_statuses
WHERE deleted = 0
AND status != 0
AND id = ?;
