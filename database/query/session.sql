-- name: GetSessionByUuid :one
SELECT * FROM sessions
WHERE uuid = ?;

-- name: CreateSession :exec
INSERT INTO sessions (uuid, user_id, expires_at, created_at, updated_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);