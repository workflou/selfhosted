-- name: CountAdmins :one
SELECT COUNT(*) FROM users
WHERE is_admin = true;

-- name: CreateAdmin :execlastid
INSERT INTO users (name, email, password, is_admin)
VALUES (?, ?, ?, true);

-- name: GetAdminByEmail :one
SELECT id, name, email, is_admin
FROM users
WHERE email = ? AND is_admin = true;

-- name: ChangeAdminPassword :exec
UPDATE users
SET password = ?
WHERE email = ? AND is_admin = true;