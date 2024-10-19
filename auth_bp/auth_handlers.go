package auth_bp

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
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

type loginUserData struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,printascii,min=8"`
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	lu := loginUserData{}
	reader := json.NewDecoder(r.Body)
	reader.Decode(&lu)

	apiCfg := apiconf.GetConfig(r)

	err := apiCfg.Validate.Struct(lu)

	if err != nil {
		validationError, ok := err.(validator.ValidationErrors)
		if !ok {
			slog.Error("error with validator definition", "error", err)
			render.RespondFailure(w, http.StatusInternalServerError, InternalServerErrorMssg)
		} else {
			render.RespondValidationFailure(w, validationError)
		}
		return
	}

	queries := database.New(apiCfg.ConnPool)
	user, err := queries.GetUserByName(r.Context(), lu.Username)
	if err != nil {
		render.RespondFailure(w, http.StatusBadRequest, "username or password is invalid")
		return
	}
	hashedPassword := saltyPassword([]byte(lu.Password), user.PasswordSalt)
	if subtle.ConstantTimeCompare(hashedPassword, user.Password) != 1 {
		render.RespondFailure(w, http.StatusBadRequest, "username or password is invalid")
		return
	}
	userData := auth.DbUserToUserData(user)
	token, err := auth.UserToToken(&userData)
	if err != nil {
		render.RespondFailure(w, http.StatusInternalServerError, TokenGenerationErrorMssg)
		return
	}
	w.Header().Set(auth.UserAuthHeader, token)
	render.RespondSuccess(w, 204, "")
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
}

type registerUserData struct {
	Username        string `json:"username" validate:"required,min=5,max=50"`
	DisplayName     string `json:"display_name" validate:"required,min=5,max=150"`
	Password        string `json:"password" validate:"required,printascii,min=8,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	ru := registerUserData{}
	reader := json.NewDecoder(r.Body)
	reader.Decode(&ru)

	apiCfg := apiconf.GetConfig(r)
	// validate incoming data
	err := apiCfg.Validate.Struct(ru)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			slog.Error("error with validator definition", "error", err)
			render.RespondFailure(w, http.StatusInternalServerError, InternalServerErrorMssg)
		} else {
			render.RespondValidationFailure(w, validationErrors)
		}
		return
	}
	// check user name with DB
	queries := database.New(apiCfg.ConnPool)
	_, err = queries.GetUserByName(r.Context(), ru.Username)
	if err == nil {
		slog.Warn("error while searching DB", "error", err)
		render.RespondFailure(w, http.StatusBadRequest, map[string]string{"username": "already exists"})
		return
	}
	// generate the password hash
	passwordSalt := make([]byte, 128)
	_, err = rand.Read(passwordSalt)
	if err != nil {
		render.RespondFailure(w, http.StatusInsufficientStorage, InsufficientStorageErrorMssg)
		return
	}
	password := saltyPassword([]byte(ru.Password), passwordSalt)
	// store in DB
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
			slog.Error(
				"could not create user",
				"error", pgErr.Message,
				"code", pgErr.Code,
				"constraint", pgErr.ConstraintName,
			)
		}
		render.RespondFailure(w, http.StatusInsufficientStorage, InsufficientStorageErrorMssg)
		return
	}
	// send back user data
	cleanUser := auth.DbUserToUserData(user)
	render.RespondSuccess(w, http.StatusCreated, cleanUser)
}

func NewRouter() *chi.Mux {
	authRouter := chi.NewRouter()

	authRouter.Post("/login", handleLogin)
	authRouter.Post("/register", handleRegister)
	authRouter.Post("/logout", handleLogout)

	return authRouter
}
