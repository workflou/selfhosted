-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: CreateUser :exec
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