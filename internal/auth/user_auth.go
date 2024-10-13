package auth

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserData struct {
	UserId     pgtype.UUID `json:"user_id"`
	UserName   string      `json:"user_name"`
	UserEmail  string      `json:"user_email"`
	CreatedAt  time.Time   `json:"created_at"`
	ExpriyTime time.Time   `json:"expiry_time"`
}

func UserToToken(u UserData) (string, error) {
	return "", nil
}

func TokenToUesr(s string) (*UserData, error) {
	return nil, nil
}
