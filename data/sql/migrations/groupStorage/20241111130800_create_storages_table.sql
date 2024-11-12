-- +goose Up
CREATE TABLE storages(
    id INTEGER PRIMARY KEY,
    group_name TEXT NOT NULL,
    storage_name TEXT NOT NULL,
    clear_slab_qty INTEGER,
    clear_block_qty INTEGER,
    cloudy_slab_qty INTEGER,
    cloudy_block_qty INTEGER,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE storages;