-- users.sql

-- name: CreateUser :one
-- add "one" user to the DB by email address
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(), -- generate a unique id
    NOW(),             -- current time
    NOW(),             -- current time
    $1                 -- gen code will input email
)
-- func generated will return these values for use in code
RETURNING *;

-- name: ResetUsers :exec
-- "reset" all users
DELETE FROM users;