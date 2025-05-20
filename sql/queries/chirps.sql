-- chirps.sql

-- name: CreateChirp :one
-- add "one" chirp to the DB, user_id is fk
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(), -- generate a unique id
    NOW(),             -- current time
    NOW(),             -- current time
    $1,                -- gen code will input body
    $2                 -- gen code will input user_id
)
-- func generated will return these values for use in code
RETURNING *;

-- name: GetChirps :many
-- select all chirps!
SELECT * FROM chirps
-- oldest to latest based on created_at
ORDER BY created_at ASC;
