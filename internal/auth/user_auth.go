package auth

import (
	"time"

	"github.com/Suryarpan/chat-api/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserData struct {
	UserId      pgtype.UUID   `json:"user_id"`
	UserName    string        `json:"user_name"`
	DisplayName string        `json:"display_name"`
	CreatedAt   time.Time     `json:"created_at"`
	ExpriyTime  time.Duration `json:"expiry_time"`
}

func DbUserToUserData(u database.User) UserData {
	return UserData{
		UserId:      u.UserID,
		UserName:    u.Username,
		DisplayName: u.DisplayName,
		CreatedAt:   u.CreatedAt,
		ExpriyTime:  time.Hour * 24,
	}
}

func UserToToken(u *UserData) (string, error) {
	return "", nil
}

func TokenToUesr(s string) (*UserData, error) {
	return nil, nil
}
