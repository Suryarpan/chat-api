package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Suryarpan/chat-api/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/pbkdf2"
)

const TokenIssuer string = "chat-api"

var (
	RegularAudience []string = []string{"user"}
	AdminAudience   []string = []string{"user", "admin"}
	secret          []byte
)

func SetupAuth() error {
	val, ok := os.LookupEnv("CHAT_API_SECRET")
	if !ok {
		return fmt.Errorf("could not find CHAT_API_SECRET")
	}
	sc, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return err
	}
	secret = sc
	return nil
}

func SaltyPassword(password, salt []byte) []byte {
	iterations := 10_000
	hashed := pbkdf2.Key(password, salt, iterations, 512, sha256.New)
	return hashed
}

type tokenData struct {
	jwt.RegisteredClaims
	UserId pgtype.UUID `json:"uid"`
}

type UserData struct {
	PvtId    int32
	UserId   pgtype.UUID
	UserName string
}

func getUserData(ctx context.Context, queries *database.Queries, t *tokenData) (UserData, error) {
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, tokenData{
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

func tokenToUser(s string, secret []byte) (*tokenData, error) {
	token, err := jwt.ParseWithClaims(
		s,
		&tokenData{},
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
	if claims, ok := token.Claims.(*tokenData); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("unknown claims type")
}
