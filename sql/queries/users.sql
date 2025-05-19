-- users.sql

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(), -- generate a unique id
    NOW(),             -- current time
    NOW(),             -- current time
    $1                 -- gen code will input email
)
RETURNING *;
