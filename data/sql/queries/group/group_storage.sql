-- name: GetStorage :one
SELECT * FROM storages
WHERE id = $1;

-- name: CreateStorage :one
INSERT INTO storages (
    group_name, storage_name, clear_slab_qty, clear_block_qty, cloudy_slab_qty, cloudy_block_qty, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: DeleteStorage :exec
DELETE FROM storages
WHERE id = $1;

-- name: UpdateStorage :exec
UPDATE storages
SET group_name = $1,
storage_name = $2,
clear_slab_qty = $3,
clear_block_qty = $4,
cloudy_slab_qty = $5,
cloudy_block_qty = $6,
created_at = $7
WHERE id = $8;

-- name: GetGroupStorages :many
SELECT * FROM storages
WHERE group_name = $1;