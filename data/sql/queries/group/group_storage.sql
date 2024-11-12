-- name: GetStorage :one
SELECT * FROM storages
WHERE id = ?;

-- name: CreateStorage :one
INSERT INTO storages (
    group_name, storage_name, clear_slab_qty, clear_block_qty, cloudy_slab_qty, cloudy_block_qty, created_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteStorage :exec
DELETE FROM storages
WHERE id = ?;

-- name: UpdateStorage :exec
UPDATE storages
SET group_name = ?,
storage_name = ?,
clear_slab_qty = ?,
clear_block_qty = ?,
cloudy_slab_qty = ?,
cloudy_block_qty = ?,
created_at = ?
WHERE id = ?;

-- name: GetGroupStorages :many
SELECT * FROM storages
WHERE group_name = ?;