package auth_bp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/Suryarpan/chat-api/internal/apiconf"
	"github.com/Suryarpan/chat-api/internal/auth"
	"github.com/Suryarpan/chat-api/internal/database"
	"github.com/Suryarpan/chat-api/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/pbkdf2"
)

func saltyPassword(password, salt []byte) []byte {
	iterations := 10_000
	hashed := pbkdf2.Key(password, salt, iterations, 512, sha256.New)
	return hashed
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
}

type registerUser struct {
	Username        string `json:"username" validate:"required,min=5,max=50"`
	DisplayName     string `json:"display_name" validate:"required,min=5,max=150"`
	Password        string `json:"password" validate:"required,printascii,min=8,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	ru := registerUser{}
	reader := json.NewDecoder(r.Body)
	reader.Decode(&ru)

	apiCfg := apiconf.GetConfig(r)
	err := apiCfg.Validate.Struct(ru)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			slog.Error("error with validator definition", "error", err)
			render.RespondFailure(w, 500, "could not process request at this time")
		} else {
			render.RespondValidationFailure(w, validationErrors)
		}
		return
	}

	queries := database.New(apiCfg.ConnPool)
	_, err = queries.GetUserByName(r.Context(), ru.Username)
	if err == nil {
		slog.Warn("error while searching DB", "error", err)
		field, _ := reflect.TypeOf(ru).FieldByName("Username")
		render.RespondFailure(w, 400, map[string]string{field.Tag.Get("json"): "already exists"})
		return
	}

	passwordSalt := make([]byte, 128)
	_, err = rand.Read(passwordSalt)
	if err != nil {
		render.RespondFailure(w, 500, "server cannot store the password at this moment")
		return
	}

	password := saltyPassword([]byte(ru.Password), passwordSalt)

	user, err := queries.CreateUser(r.Context(), database.CreateUserParams{
		Username:     ru.Username,
		DisplayName:  ru.DisplayName,
		Password:     password,
		PasswordSalt: passwordSalt,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			slog.Error("could not create user", "error", pgErr.Message, "code", pgErr.Code, "constraint", pgErr.ConstraintName)
		}
		render.RespondFailure(w, 500, "somthing went wrong while creating user")
		return
	}

	cleanUser := auth.DbUserToUserData(user)
	render.RespondSuccess(w, 201, cleanUser)
}

func NewRouter() *chi.Mux {
	authRouter := chi.NewRouter()

	authRouter.Post("/login", handleLogin)
	authRouter.Post("/register", handleRegister)
	authRouter.Post("/logout", handleLogout)

	return authRouter
}
