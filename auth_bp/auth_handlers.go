package auth_bp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
}

func NewRouter() *chi.Mux {
	authRouter := chi.NewRouter()

	authRouter.Post("/login", handleLogin)
	authRouter.Post("/logout", handleLogout)

	return authRouter
}
