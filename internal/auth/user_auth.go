package auth

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Suryarpan/chat-api/internal/database"
	"github.com/Suryarpan/chat-api/render"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ctxKeyUserData string

const (
	ctxUserDataKey ctxKeyUserData = "CHAT_API_USER_DATA"
	UserAuthHeader string         = "Authorization"
	TokenPrefix    string         = "Bearer"
)

var (
	Secret []byte
)

func init() {
	val, ok := os.LookupEnv("CHAT_API_SECRET")
	if !ok {
		panic("Error: could not find CHAT_API_SECRET")
	}
	secret, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		panic(err)
	}
	Secret = secret
}

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
		ExpriyTime:  time.Hour * 24 / time.Second,
	}
}

func UserToToken(u *UserData) (string, error) {
	// jwt.NewWithClaims()
	return "", nil
}

func TokenToUser(s string, secret []byte) (*UserData, error) {
	return nil, nil
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(UserAuthHeader)
		parts := strings.Split(token, " ")
		mssg := "please authenticate before proceeding"
		if len(parts) != 2 {
			render.RespondFailure(w, http.StatusUnauthorized, mssg)
			return
		} else if parts[0] != TokenPrefix {
			render.RespondFailure(w, http.StatusUnauthorized, mssg)
			return
		}
		data, err := TokenToUser(parts[1], Secret)
		if err != nil {
			render.RespondFailure(w, http.StatusUnauthorized, mssg)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserDataKey, data)
		rr := r.WithContext(ctx)
		next.ServeHTTP(w, rr)
	})
}
