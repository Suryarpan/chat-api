package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Suryarpan/chat-api/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ctxKeyUserData string

const (
	ctxUserDataKey ctxKeyUserData = "CHAT_API_USER_DATA"
	UserAuthHeader string         = "Authorization"
	TokenPrefix    string         = "Bearer"
	TokenIssuer    string         = "chat-api"
	LenPrefix      int            = len(TokenPrefix + " ")
)

var (
	RegularAudience []string = []string{"user"}
	AdminAudience   []string = []string{"user", "admin"}
	secret          []byte
)

func init() {
	val, ok := os.LookupEnv("CHAT_API_SECRET")
	if !ok {
		panic("Error: could not find CHAT_API_SECRET")
	}
	sc, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		panic(err)
	}
	secret = sc
}

type TokenData struct {
	jwt.RegisteredClaims
	UserId pgtype.UUID `json:"uid"`
}

type UserData struct {
	PvtId    int32
	UserId   pgtype.UUID
	UserName string
}

func GetUserData(ctx context.Context, queries database.Queries, t TokenData) (UserData, error) {
	user, err := queries.GetUserByNameAndUuid(ctx, database.GetUserByNameAndUuidParams{
		UserID:   t.UserId,
		Username: t.Subject,
	})
	if err != nil {
		return UserData{}, err
	}
	return UserData{
		PvtId:    user.PvtID,
		UserId:   user.UserID,
		UserName: user.Username,
	}, nil
}

func UserToToken(u database.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, TokenData{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  RegularAudience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24).UTC()),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    TokenIssuer,
			Subject:   u.Username,
		},
		UserId: u.UserID,
	})
	return token.SignedString(secret)
}

func TokenToUser(s string, secret []byte) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(
		s,
		&TokenData{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		},
		jwt.WithAudience("user"),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer("chat-api"),
		jwt.WithLeeway(time.Second*10),
		jwt.WithJSONNumber(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS384.Name}),
	)
	if err != nil {
		slog.Warn("invalid token encountered", "error", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*TokenData); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("unknown claims type")
}
