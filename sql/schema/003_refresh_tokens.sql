-- 003_refresh_tokens.sql
-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,      -- unique note string
    created_at TIMESTAMP NOT NULL, -- for auditing
    updated_at TIMESTAMP NOT NULL, -- for auditing
    user_id UUID NOT NULL,         -- user id for fk
    expires_at TIMESTAMP NOT NULL, -- expiration checking
    revoked_at TIMESTAMP NULL,      -- defaults to "null"
    -- link user_id to refresh_token as fk
    FOREIGN KEY (user_id) -- select fk
        REFERENCES users (id) -- match with id in users
        ON DELETE CASCADE -- prevents orphan refresh_tokens
);

-- +goose Down
DROP TABLE refresh_tokens;