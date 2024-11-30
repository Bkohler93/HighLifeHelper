-- +goose Up
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,            -- SERIAL for auto-incrementing ID
    uuid UUID NOT NULL,               -- Use the UUID type for universally unique identifiers
    login_id TEXT NOT NULL,           -- TEXT for login identifiers
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(), -- Default to current timestamp
    expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL               -- Explicitly define expiration timestamp
);

-- +goose Down
DROP TABLE sessions;