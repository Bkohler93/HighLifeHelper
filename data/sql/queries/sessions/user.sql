-- name: GetUser :one
SELECT * FROM users
WHERE login_id = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    login_id, pw_hash, created_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE login_id = $1;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1;