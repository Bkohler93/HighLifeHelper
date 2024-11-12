-- name: GetSession :one
SELECT * FROM sessions
WHERE uuid = ? LIMIT 1;

-- name: GetSessionById :one
SELECT * FROM sessions
WHERE id = ? LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (
    uuid, login_id, created_at, expires_at
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;