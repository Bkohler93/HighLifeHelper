-- +goose Up
CREATE TABLE sessions(
    id INTEGER PRIMARY KEY,
    uuid TEXT NOT NULL,
    login_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE sessions;