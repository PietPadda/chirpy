-- 002_chirps.sql
-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY, -- our pk
    created_at TIMESTAMP NOT NULL, -- audit trail
    updated_at TIMESTAMP NOT NULL, -- audit trail
    body TEXT NOT NULL, -- da chirp!
    user_id UUID NOT NULL, -- our fk
    -- link user_id to chirp as fk
    FOREIGN KEY (user_id) -- select fk
        REFERENCES users (id) -- match with id in users
        ON DELETE CASCADE -- prevents orphan chirps
);

-- +goose Down
DROP TABLE chirps;
