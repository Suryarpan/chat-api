package apiconf

import (
	"context"
	"log/slog"
	"net/http"
)

type ctxKeyApiConfig string
const chatApiConfigKey ctxKeyApiConfig = "CHAT_API_DB_URL"

type ApiConfig struct {
	DBUrl string
}

func ApiConfigure(apiCfg *ApiConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), chatApiConfigKey, *apiCfg)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetConfig(r *http.Request) ApiConfig {
	apiCfg, ok := r.Context().Value(chatApiConfigKey).(ApiConfig)
	if !ok {
		slog.Error("ChatApiConfigKey was overwrittten", "config", apiCfg)
		panic("cannot prceed further with corrupted config")
	}
	return apiCfg
}
