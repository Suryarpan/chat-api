-- name: CreateUser :one
INSERT INTO users (
    id, username, display_name, password, password_salt, created_at, updated_at, last_logged_in
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, NULL
) RETURNING *;
 