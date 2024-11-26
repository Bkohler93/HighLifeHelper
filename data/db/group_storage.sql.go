// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: group_storage.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createStorage = `-- name: CreateStorage :one
INSERT INTO storages (
    group_name, storage_name, clear_slab_qty, clear_block_qty, cloudy_slab_qty, cloudy_block_qty, created_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
)
RETURNING id, group_name, storage_name, clear_slab_qty, clear_block_qty, cloudy_slab_qty, cloudy_block_qty, created_at
`

type CreateStorageParams struct {
	GroupName      string
	StorageName    string
	ClearSlabQty   sql.NullInt64
	ClearBlockQty  sql.NullInt64
	CloudySlabQty  sql.NullInt64
	CloudyBlockQty sql.NullInt64
	CreatedAt      time.Time
}

func (q *Queries) CreateStorage(ctx context.Context, arg CreateStorageParams) (Storage, error) {
	row := q.db.QueryRowContext(ctx, createStorage,
		arg.GroupName,
		arg.StorageName,
		arg.ClearSlabQty,
		arg.ClearBlockQty,
		arg.CloudySlabQty,
		arg.CloudyBlockQty,
		arg.CreatedAt,
	)
	var i Storage
	err := row.Scan(
		&i.ID,
		&i.GroupName,
		&i.StorageName,
		&i.ClearSlabQty,
		&i.ClearBlockQty,
		&i.CloudySlabQty,
		&i.CloudyBlockQty,
		&i.CreatedAt,
	)
	return i, err
}

const deleteStorage = `-- name: DeleteStorage :exec
DELETE FROM storages
WHERE id = ?
`

func (q *Queries) DeleteStorage(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteStorage, id)
	return err
}

const getGroupStorages = `-- name: GetGroupStorages :many
SELECT id, group_name, storage_name, clear_slab_qty, clear_block_qty, cloudy_slab_qty, cloudy_block_qty, created_at FROM storages
WHERE group_name = ?
`

func (q *Queries) GetGroupStorages(ctx context.Context, groupName string) ([]Storage, error) {
	rows, err := q.db.QueryContext(ctx, getGroupStorages, groupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Storage
	for rows.Next() {
		var i Storage
		if err := rows.Scan(
			&i.ID,
			&i.GroupName,
			&i.StorageName,
			&i.ClearSlabQty,
			&i.ClearBlockQty,
			&i.CloudySlabQty,
			&i.CloudyBlockQty,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getStorage = `-- name: GetStorage :one
SELECT id, group_name, storage_name, clear_slab_qty, clear_block_qty, cloudy_slab_qty, cloudy_block_qty, created_at FROM storages
WHERE id = ?
`

func (q *Queries) GetStorage(ctx context.Context, id int64) (Storage, error) {
	row := q.db.QueryRowContext(ctx, getStorage, id)
	var i Storage
	err := row.Scan(
		&i.ID,
		&i.GroupName,
		&i.StorageName,
		&i.ClearSlabQty,
		&i.ClearBlockQty,
		&i.CloudySlabQty,
		&i.CloudyBlockQty,
		&i.CreatedAt,
	)
	return i, err
}

const updateStorage = `-- name: UpdateStorage :exec
UPDATE storages
SET group_name = ?,
storage_name = ?,
clear_slab_qty = ?,
clear_block_qty = ?,
cloudy_slab_qty = ?,
cloudy_block_qty = ?,
created_at = ?
WHERE id = ?
`

type UpdateStorageParams struct {
	GroupName      string
	StorageName    string
	ClearSlabQty   sql.NullInt64
	ClearBlockQty  sql.NullInt64
	CloudySlabQty  sql.NullInt64
	CloudyBlockQty sql.NullInt64
	CreatedAt      time.Time
	ID             int64
}

func (q *Queries) UpdateStorage(ctx context.Context, arg UpdateStorageParams) error {
	_, err := q.db.ExecContext(ctx, updateStorage,
		arg.GroupName,
		arg.StorageName,
		arg.ClearSlabQty,
		arg.ClearBlockQty,
		arg.CloudySlabQty,
		arg.CloudyBlockQty,
		arg.CreatedAt,
		arg.ID,
	)
	return err
}
