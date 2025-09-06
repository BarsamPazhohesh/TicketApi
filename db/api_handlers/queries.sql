-- name: AddApiHandler :one
INSERT INTO api_handlers (handler, method, description) VALUES (?, ?, ?) RETURNING id;