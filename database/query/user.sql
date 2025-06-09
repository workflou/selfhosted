-- name: GetUserByEmail :one
SELECT users.*, sqlc.embed(teams) FROM users
LEFT JOIN teams ON users.current_team_id = teams.id
WHERE users.email = ?;

-- name: CreateUser :execlastid
INSERT INTO users (name, email, password, created_at, updated_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: UpdateUserName :exec
UPDATE users
SET name = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateUserAvatar :exec
UPDATE users
SET avatar = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;