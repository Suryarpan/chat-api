-- name: CreateUser :one
INSERT INTO users (
    user_id, username, display_name, password, password_salt, created_at, updated_at, last_logged_in
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, NULL
) RETURNING *;

-- name: GetUserByUuid :one
SELECT *
FROM users
WHERE user_id = $1;

-- name: GetUserByName :one
SELECT *
FROM users
WHERE username = $1;

-- name: GetUserByNameAndUuid :one
SELECT *
FROM users
WHERE user_id = $1 AND username = $2;

-- name: GetUserById :one
SELECT *
FROM users
WHERE pvt_id = $1;

-- name: UpdateLoggedInTime :exec
UPDATE users
SET last_logged_in = $1
WHERE pvt_id = $2;

-- name: UpdateUserDetails :one
UPDATE users
SET username = $1, display_name = $2, password = $3, updated_at = $4
WHERE pvt_id = $5
RETURNING *;

-- name: DeleteUserDetails :one
DELETE
FROM users
WHERE pvt_id = $1
RETURNING *;