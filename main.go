package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Suryarpan/chat-api/internal/apiconf"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Config struct {
	DBUrl string
}

func setUpMiddlewares(r *chi.Mux) error {
	if r == nil {
		return errors.New("please provide and router")
	}
	r.Use(apiconf.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.ContentCharset("UTF-8", "Latin-1"))
	r.Use(middleware.AllowContentType("application/json", "text/xml"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	return nil
}

func main() {
	logger := slog.New(*apiconf.LoggerConfig())
	slog.SetDefault(logger)

	_ = Config{
		DBUrl: apiconf.DBUrlConfig(),
	}

	// router setup
	mainRouter := chi.NewRouter()
	err := setUpMiddlewares(mainRouter)
	if err != nil {
		panic(fmt.Sprintf("Error: could not setup middlewares: %v", err))
	}
	// base api setup
	apiV1router := chi.NewRouter()
	apiV1router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	// auth api setup
	// user api setup
	// chat data setup
	// admin setup
	// router run
	mainRouter.Mount("/api/v1", apiV1router)
	port, ok := os.LookupEnv("CHAT_API_PORT")
	if !ok {
		panic("Error: could not find CHAT_API_PORT environment variable")
	}
	logger.Info("starting router", "port", port)

	server := &http.Server{
		Handler: mainRouter,
		Addr:    ":" + port,
	}
	err = server.ListenAndServe()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
	}
}
