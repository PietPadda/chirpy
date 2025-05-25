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

-- name: UpdateUserLogin :one 
UPDATE users 
SET 
  updated_at = NOW(),  -- audit trail
  email = $2, -- user provides new email
  hashed_password = $3 -- user provides new password
WHERE id = $1 -- use userid from token get bearer (unique as it's a pk) 
-- return only updated_at to match resp timestamp (rest are inputs from code, no need to return)
RETURNING updated_at;

-- name: GetUserByID :one
-- select one user by user_id
SELECT * FROM users
-- by user id as input
WHERE id = $1
LIMIT 1;
