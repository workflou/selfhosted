-- name: CreateTeam :execlastid
INSERT INTO teams (uuid, name, created_at, updated_at)
VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: AddMemberToTeam :exec
INSERT INTO members (team_id, user_id, role, created_at, updated_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: GetUserTeams :many
SELECT *, sqlc.embed(teams), sqlc.embed(users)
FROM members
LEFT JOIN teams ON members.team_id = teams.id
LEFT JOIN users ON members.user_id = users.id
WHERE members.user_id = ?;

-- name: SetCurrentTeam :exec
UPDATE users
SET current_team_id = ?
WHERE id = ?;