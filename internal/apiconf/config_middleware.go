package apiconf

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKeyApiConfig string

const chatApiConfigKey ctxKeyApiConfig = "CHAT_API_DB_URL"

type ApiConfig struct {
	ConnPool *pgxpool.Pool
	Validate *validator.Validate
}

func SetupPool() (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(DBUrlConfig())
	if err != nil {
		return nil, err
	}
	dbConfig.MaxConns = 10
	dbConfig.MinConns = 0
	dbConfig.MaxConnLifetimeJitter = time.Hour * 1
	dbConfig.MaxConnIdleTime = time.Minute * 5
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 10

	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}
	err = connPool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return connPool, nil
}

func ApiConfigure(connPool *pgxpool.Pool) func(http.Handler) http.Handler {
	apiCfg := ApiConfig{
		ConnPool: connPool,
		Validate: validator.New(validator.WithRequiredStructEnabled()),
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), chatApiConfigKey, apiCfg)
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
