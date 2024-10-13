// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id, username, display_name, password, password_salt, created_at, updated_at, last_logged_in
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NULL
) RETURNING id, username, display_name, password, password_salt, created_at, updated_at, last_logged_in
`

type CreateUserParams struct {
	ID           pgtype.UUID      `json:"id"`
	Username     string           `json:"username"`
	DisplayName  string           `json:"display_name"`
	Password     []byte           `json:"password"`
	PasswordSalt []byte           `json:"password_salt"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
	UpdatedAt    pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.DisplayName,
		arg.Password,
		arg.PasswordSalt,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.DisplayName,
		&i.Password,
		&i.PasswordSalt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLoggedIn,
	)
	return i, err
}
