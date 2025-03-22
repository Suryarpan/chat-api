package main

import (
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
	insufficientStorageMessageError = "could not create message at this moment"
)

type createMessageData struct {
	ToUserId     pgtype.UUID `json:"to_user_id" validate:"required"`
	MssgType     string      `json:"mssg_type" validate:"required,oneof=normal reply reaction"`
	AttachMssgId int64       `json:"attach_mssg_id" validate:"omitempty,min=1"`
	MssgBody     string      `json:"mssg_body" validate:"required,min=1,printascii|alphanumunicode"`
}

func handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	data := createMessageData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		slog.Warn("could not decode incoming data")
		render.RespondFailure(w, 400, "could not decode data")
		return
	}

	apiCfg := apiconf.GetConfig(r)
	// validate incoming data
	slog.Info("validating user data", "data", data)
	err = apiCfg.Validate.Struct(data)
	if err != nil {
		slog.Warn("could not validate data")
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			slog.Error("error with validator definition", "error", err)
			render.RespondFailure(w, http.StatusInternalServerError, internalServerErrorMssg)
		} else {
			render.RespondValidationFailure(w, validationErrors)
		}
		return
	}
	fromUser := auth.GetUserData(r)

	slog.Info("fetching to send user data")
	queries := database.New(apiCfg.ConnPool)
	toUser, err := queries.GetUserByUuid(r.Context(), data.ToUserId)
	if err != nil {
		render.RespondFailure(w, http.StatusBadRequest, "could not find user to send to")
		return
	}
	slog.Debug("received user data", "to user", toUser)

	c, err := apiCfg.ConnPool.Acquire(r.Context())
	defer c.Release()
	if err != nil {
		render.RespondFailure(w, http.StatusInternalServerError, internalServerErrorMssg)
		return
	}

	slog.Info("starting transaction")
	tx, err := c.Begin(r.Context())
	defer tx.Rollback(r.Context())
	if err != nil {
		render.RespondFailure(w, http.StatusInternalServerError, internalServerErrorMssg)
		return
	}

	slog.Debug("creating message meta entry")
	txQuery := queries.WithTx(tx)
	mssgMeta, err := txQuery.CreateMessage(r.Context(), database.CreateMessageParams{
		FromPvtID:  fromUser.PvtID,
		ToPvtID:    toUser.PvtID,
		MssgStatus: database.MessageStatusSent,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})
	if err != nil {
		render.RespondFailure(w, http.StatusInsufficientStorage, insufficientStorageMessageError)
		return
	}

	slog.Debug("creating type information entry of message", "mssg id", mssgMeta.MssgID)
	_, err = txQuery.CreateMessageType(r.Context(), database.CreateMessageTypeParams{
		MssgID:   mssgMeta.MssgID,
		MssgType: database.MessageType(data.MssgType),
		AttachMssgID: pgtype.Int8{
			Int64: data.AttachMssgId,
			Valid: data.AttachMssgId != 0,
		},
	})
	if err != nil {
		render.RespondFailure(w, http.StatusInsufficientStorage, insufficientStorageMessageError)
		return
	}

	slog.Debug("creating body entry of message", "mssg id", mssgMeta.MssgID)
	_, err = txQuery.CreateMessageText(r.Context(), database.CreateMessageTextParams{
		MssgID:   mssgMeta.MssgID,
		MssgBody: data.MssgBody,
	})
	if err != nil {
		render.RespondFailure(w, http.StatusInsufficientStorage, insufficientStorageMessageError)
		return
	}

	slog.Debug("fetching public data from database", "mssg_id", mssgMeta.MssgID)
	mssgContent, err := txQuery.GetMessageByIdPublic(r.Context(), mssgMeta.MssgID)
	if err != nil {
		render.RespondFailure(w, http.StatusInsufficientStorage, insufficientStorageMessageError)
		return
	}

	slog.Debug("commiting the db writes", "mssg_id", mssgMeta.MssgID)
	err = tx.Commit(r.Context())
	if err != nil {
		render.RespondFailure(w, http.StatusInsufficientStorage, insufficientStorageMessageError)
		return
	}

	slog.Info("sending back reponse", "message", mssgContent)
	render.RespondSuccess(w, http.StatusOK, mssgContent)
}






  
func MessageRouter() *chi.Mux {
	router := chi.NewMux()

	router.Post("/", handleCreateMessage)

	return router
}
