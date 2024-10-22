package apiconf

import (
	"fmt"
	"log/slog"
	"log/syslog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggerConfig() *slog.Handler {
	level, ok := os.LookupEnv("CHAT_API_CONFIG")
	if !ok {
		level = "prod"
	}
	var logger slog.Handler

	switch level {
	case "dev":
		logger = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	case "prod":
		logWriter, err := syslog.New(syslog.LOG_NOTICE, "chat-api")
		if err != nil {
			panic("Error: could not setup connection to syslog")
		}
		logger = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
	default:
		panic(fmt.Sprintf("Error: could not understand log level: %s\n", level))
	}
	return &logger
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		next.ServeHTTP(ww, r)

		slog.Info(
			fmt.Sprintf("%s://%s%s %s", scheme, r.Host, r.RequestURI, r.Proto),
			"from", r.RemoteAddr,
			"method", r.Method,
			"status", ww.Status(),
			"length", ww.BytesWritten(),
			"time_taken", time.Since(t1),
		)
	})
}
