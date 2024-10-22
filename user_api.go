package main

import (
	"crypto/rand"
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
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	insufficientStorageErrorMssg = "could not create user at this moment"
)

type PublicUserDetails struct {
	UserID       pgtype.UUID      `json:"user_id"`
	Username     string           `json:"username"`
	DisplayName  string           `json:"display_name"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	LastLoggedIn pgtype.Timestamp `json:"last_logged_in"`
}

func handleGetUserDetail(w http.ResponseWriter, r *http.Request) {
	apiCfg := apiconf.GetConfig(r)
	userData := auth.GetUserData(r)

	queries := database.New(apiCfg.ConnPool)
	user, err := queries.GetUserById(r.Context(), userData.PvtId)
	if err != nil {
		render.RespondFailure(w, http.StatusInternalServerError, "could not fetch user data")
		return
	}

	publicData := PublicUserDetails{
		UserID:       user.UserID,
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		LastLoggedIn: user.LastLoggedIn,
	}
	render.RespondSuccess(w, http.StatusOK, publicData)
}

type createUserData struct {
	Username    string `json:"username" validate:"required,min=5,max=50"`
	DisplayName string `json:"display_name" validate:"required,min=5,max=150"`
	Password    string `json:"password" validate:"required,printascii,min=8"`
}

func handleUserCreate(w http.ResponseWriter, r *http.Request) {
	ru := createUserData{}
	reader := json.NewDecoder(r.Body)
	reader.Decode(&ru)

	apiCfg := apiconf.GetConfig(r)
	// validate incoming data
	err := apiCfg.Validate.Struct(ru)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			slog.Error("error with validator definition", "error", err)
			render.RespondFailure(w, http.StatusInternalServerError, internalServerErrorMssg)
		} else {
			render.RespondValidationFailure(w, validationErrors)
		}
		return
	}
	// check user name with DB
	queries := database.New(apiCfg.ConnPool)
	_, err = queries.GetUserByName(r.Context(), ru.Username)
	if err == nil {
		render.RespondFailure(w, http.StatusNotAcceptable, map[string]string{"username": "already exists"})
		return
	}
	// generate the password hash
	passwordSalt := make([]byte, 128)
	_, err = rand.Read(passwordSalt)
	if err != nil {
		render.RespondFailure(w, http.StatusInsufficientStorage, insufficientStorageErrorMssg)
		return
	}
	password := auth.SaltyPassword([]byte(ru.Password), passwordSalt)
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
		render.RespondFailure(w, http.StatusInsufficientStorage, insufficientStorageErrorMssg)
		return
	}
	// send back user data
	render.RespondSuccess(w, http.StatusCreated, PublicUserDetails{
		UserID:      user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	})
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {

}

func UserRouter() *chi.Mux {
	router := chi.NewMux()

	router.With(auth.Authentication).Get("/", handleGetUserDetail)
	router.Post("/", handleUserCreate)
	router.Patch("/", handleUpdateUser)

	return router
}
