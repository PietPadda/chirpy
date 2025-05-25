-- refresh_tokens.sql

-- name: CreateRefreshToken :one
-- add "one" refresh token to the DB, user_id is fk
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,     -- insert token string
    NOW(),  -- current time
    NOW(),  -- current time
    $2,     -- insert user id fk
    $3,     -- insert expiration time
    NULL    -- default to null
)
-- func generated will return these values for use in code
RETURNING *;

-- name: GetUserFromRefreshToken :one
-- select one user from refresh token
SELECT user_id, expires_at, revoked_at FROM refresh_tokens
-- by refresh token input
WHERE token = $1
LIMIT 1;

-- name: RevokeRefreshToken :exec 
UPDATE refresh_tokens 
SET 
  updated_at = NOW(), 
  revoked_at = NOW() 
WHERE token = $1; -- use token string (unique as it's a pk) 