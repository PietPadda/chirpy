-- 004_users_is_chirpy_red.sql
-- +goose Up
ALTER TABLE users
-- premium col added, true/false, default to false
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE
;

-- +goose Down
ALTER TABLE users
-- drop the col to undo
DROP COLUMN is_chirpy_red;