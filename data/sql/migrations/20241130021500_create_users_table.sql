-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,                
    login_id TEXT NOT NULL UNIQUE,        
    pw_hash TEXT NOT NULL,                
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW() 
);

-- +goose Down
DROP TABLE users;