package main

import (
	"crypto/subtle"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Suryarpan/chat-api/internal/apiconf"
	"github.com/Suryarpan/chat-api/internal/auth"
	"github.com/Suryarpan/chat-api/internal/database"
	"github.com/Suryarpan/chat-api/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	internalServerErrorMssg  = "could not process request at this time"
	tokenGenerationErrorMssg = "could not login user at this time"
)

type loginUserData struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,printascii,min=8"`
}

type loginResponse struct {
	Token        string           `json:"token"`
	TokenType    string           `json:"token_type"`
	Username     string           `json:"username"`
	DisplayName  string           `json:"display_name"`
	LastLoggedIn pgtype.Timestamp `json:"last_logged_in"`
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
			render.RespondFailure(w, http.StatusInternalServerError, internalServerErrorMssg)
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
	hashedPassword := auth.SaltyPassword([]byte(lu.Password), user.PasswordSalt)
	if subtle.ConstantTimeCompare(hashedPassword, user.Password) != 1 {
		render.RespondFailure(w, http.StatusBadRequest, "username or password is invalid")
		return
	}

	err = queries.UpdateLoggedInTime(r.Context(), database.UpdateLoggedInTimeParams{
		LastLoggedIn: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		PvtID: user.PvtID,
	})
	if err != nil {
		render.RespondFailure(w, http.StatusInsufficientStorage, tokenGenerationErrorMssg)
	}

	token, err := auth.UserToToken(user)
	if err != nil {
		render.RespondFailure(w, http.StatusInternalServerError, tokenGenerationErrorMssg)
		return
	}
	render.RespondSuccess(w, http.StatusOK, loginResponse{
		Token:        token,
		TokenType:    auth.TokenPrefix,
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		LastLoggedIn: user.LastLoggedIn,
	})
}

func AuthRouter() *chi.Mux {
	authRouter := chi.NewRouter()

	authRouter.Post("/login", handleLogin)

	return authRouter
}
