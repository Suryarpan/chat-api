package auth_bp

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Suryarpan/chat-api/internal/apiconf"
	"github.com/Suryarpan/chat-api/internal/database"
	"github.com/Suryarpan/chat-api/render"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
}

type registerUser struct {
	Username        string `json:"username"`
	DisplayName     string `json:"display_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	ru := registerUser{}
	reader := json.NewDecoder(r.Body)
	reader.Decode(&ru)

	// perform checks
	if ru.Password != ru.ConfirmPassword {
		render.RepondFailure(w, 400, "provided password does not match confirm password")
		return
	}

	apiCfg := apiconf.GetConfig(r)
	conn, err := pgx.Connect(r.Context(), apiCfg.DBUrl)
	defer func() {
		err = conn.Close(r.Context())
		if err != nil {
			slog.Error("error while closing DB connection", "error", err)
		}
	}()

	if err != nil {
		render.RepondFailure(w, 500, "server cannot access user data")
		return
	}

	queries := database.New(conn)
	queries.CreateUser(r.Context(), database.CreateUserParams{
		Username: ru.Username,

	})
}

func NewRouter() *chi.Mux {
	authRouter := chi.NewRouter()

	authRouter.Post("/login", handleLogin)
	authRouter.Post("/register", handleRegister)
	authRouter.Post("/logout", handleLogout)

	return authRouter
}
