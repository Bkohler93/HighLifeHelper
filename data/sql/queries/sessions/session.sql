-- name: GetSession :one
SELECT * FROM sessions
WHERE uuid = $1 LIMIT 1;

-- name: GetSessionById :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (
    uuid, login_id, created_at, expires_at
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;