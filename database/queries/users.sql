-- name: CreateUser :one
INSERT INTO users (
    id, name, email, password, password_salt, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;
 