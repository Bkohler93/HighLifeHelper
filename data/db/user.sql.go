// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    login_id, pw_hash, created_at
) VALUES (
    $1, $2, $3
)
RETURNING id, login_id, pw_hash, created_at
`

type CreateUserParams struct {
	LoginID   string
	PwHash    string
	CreatedAt time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.LoginID, arg.PwHash, arg.CreatedAt)
	var i User
	err := row.Scan(
		&i.ID,
		&i.LoginID,
		&i.PwHash,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE login_id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, loginID string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, loginID)
	return err
}

const deleteUserById = `-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUserById(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUserById, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, login_id, pw_hash, created_at FROM users
WHERE login_id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, loginID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, loginID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.LoginID,
		&i.PwHash,
		&i.CreatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, login_id, pw_hash, created_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.LoginID,
		&i.PwHash,
		&i.CreatedAt,
	)
	return i, err
}
