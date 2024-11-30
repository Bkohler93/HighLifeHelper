-- +goose Up
CREATE TABLE storages (
    id SERIAL PRIMARY KEY,
    group_name TEXT NOT NULL,
    storage_name TEXT NOT NULL,
    clear_slab_qty INTEGER DEFAULT 0,
    clear_block_qty INTEGER DEFAULT 0,
    cloudy_slab_qty INTEGER DEFAULT 0,
    cloudy_block_qty INTEGER DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW() -- Use a default value
);

-- +goose Down
DROP TABLE storages;