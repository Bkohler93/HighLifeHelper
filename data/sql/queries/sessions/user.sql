-- name: GetUser :one
SELECT * FROM users
WHERE login_id = ? LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    login_id, pw_hash, created_at
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE login_id = ?;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = ?;