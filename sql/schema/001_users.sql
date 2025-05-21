-- 001_users.sql
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY, -- unique user id
    created_at TIMESTAMP NOT NULL, -- for auditing
    updated_at TIMESTAMP NOT NULL, -- for auditing
    email TEXT NOT NULL UNIQUE, -- user login email
    hashed_password TEXT NOT NULL DEFAULT 'unset' -- defaults to "unset"
);

-- +goose Down
DROP TABLE users;