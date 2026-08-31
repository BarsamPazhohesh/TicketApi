-- name: AddTicketStatus :exec
INSERT INTO ticket_statuses (title, description) VALUES (?, ?);

-- name: GetAllActiveTicketStatuses :many
SELECT * FROM ticket_statuses
WHERE deleted_at IS NULL
AND status != 0;

-- name: GetActiveTicketStatusById :one
SELECT * FROM ticket_statuses
WHERE deleted_at IS NULL
AND status != 0
AND id = ?;
