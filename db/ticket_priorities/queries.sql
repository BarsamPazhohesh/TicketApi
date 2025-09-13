-- name: AddTicketPriority :exec
INSERT INTO ticket_priorities (user_id, ticket_type_id, priority) VALUES (?, ?, ?);

-- name: GetTicketPriorityByID :one
SELECT * FROM ticket_priorities WHERE   deleted =0
