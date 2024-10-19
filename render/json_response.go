package render

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type errorMsg struct {
	Message interface{} `json:"message"`
}


func RespondValidationFailure(w http.ResponseWriter, validationErrors validator.ValidationErrors) {
	errorMssgs := make(map[string]string)
	for _, fieldError := range validationErrors {
		errorMssgs[fieldError.Field()] = fmt.Sprintf("failed on %s with value '%s'", fieldError.ActualTag(), fieldError.Value())
	}
	RespondFailure(w, 400, errorMssgs)
}

func RespondFailure(w http.ResponseWriter, code int, msg interface{}) {
	if 399 < code && code < 499 {
		slog.Warn("bad user data received", "error", msg)
	} else if 499 < code && code < 599 {
		slog.Error("something went wrong while responding", "error", msg)
	} else {
		slog.Error("please provide valid error code", "code", code, "error", msg)
	}

	err := errorMsg{
		Message: msg,
	}

	RespondSuccess(w, code, err)
}

func RespondSuccess(w http.ResponseWriter, code int, obj interface{}) {
	w.Header().Set("content-type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		w.Write([]byte(""))
		return
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(obj)
	if err != nil {
		slog.Error("something went wrong", "error", err)
		panic("error while encoding json")
	}
}
