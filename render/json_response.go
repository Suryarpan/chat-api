package render

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorMsg struct {
	Message string `json:"message"`
}

func RepondFailure(w http.ResponseWriter, code int, msg string) {
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
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	err := enc.Encode(obj)
	if err != nil {
		slog.Error("something went wrong", "error", err)
		panic("error while encoding json")
	}
}
