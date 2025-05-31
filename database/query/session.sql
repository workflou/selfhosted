-- name: GetSessionByUuid :one
SELECT sessions.*, sqlc.embed(users) FROM sessions
LEFT JOIN users ON sessions.user_id = users.id
WHERE sessions.uuid = ?;

-- name: CreateSession :exec
INSERT INTO sessions (uuid, user_id, expires_at, created_at, updated_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE uuid = ?;