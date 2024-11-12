-- +goose Up
CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    login_id TEXT NOT NULL,
    pw_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;