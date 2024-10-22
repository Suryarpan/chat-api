package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Suryarpan/chat-api/internal/apiconf"
	"github.com/Suryarpan/chat-api/internal/database"
	"github.com/Suryarpan/chat-api/render"
	"github.com/jackc/pgx/v5"
)

type ctxKeyUserData string

const (
	ctxUserDataKey ctxKeyUserData = "CHAT_API_USER_DATA"
	UserAuthHeader string         = "Authorization"
	TokenPrefix    string         = "Bearer"
	LenPrefix      int            = len(TokenPrefix + " ")
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(UserAuthHeader)
		mssg := "please authenticate before proceeding"
		if token == "" {
			render.RespondFailure(w, http.StatusUnauthorized, mssg)
			return
		} else if !strings.HasPrefix(token, "Bearer ") {
			render.RespondFailure(w, http.StatusUnauthorized, mssg)
			return
		}
		data, err := tokenToUser(token[LenPrefix:], secret)
		if err != nil {
			render.RespondFailure(w, http.StatusUnauthorized, mssg)
			return
		}
		apiCfg := apiconf.GetConfig(r)
		queries := database.New(apiCfg.ConnPool)
		user, err := getUserData(r.Context(), queries, data)
		if err == pgx.ErrNoRows {
			render.RespondFailure(w, http.StatusUnauthorized, mssg)
			return
		} else if err != nil {
			render.RespondFailure(w, http.StatusInternalServerError, "could not login at this time")
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserDataKey, user)
		rr := r.WithContext(ctx)
		next.ServeHTTP(w, rr)
	})
}

func GetUserData(r *http.Request) database.User {
	data, ok := r.Context().Value(ctxUserDataKey).(database.User)
	if !ok {
		slog.Error("UserDataKey was overwrittten", "user", data)
		panic("cannot prceed further with corrupted user data")
	}
	return data
}
