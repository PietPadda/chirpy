-- users.sql

-- name: CreateUser :one
-- add "one" user to the DB by email address
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), -- generate a unique id
    NOW(),             -- current time
    NOW(),             -- current time
    $1,                -- gen code will input email
    $2                 -- insert hashed pw via handler
)
-- func generated will return these values for use in code
RETURNING *;

-- name: ResetUsers :exec
-- "reset" all users
DELETE FROM users;

-- name: GetUserByEmail :one
-- select one user by email
SELECT * FROM users
-- by user email as input
WHERE email = $1
LIMIT 1;