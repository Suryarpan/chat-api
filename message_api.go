package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleCreateMessage(w http.ResponseWriter, r *http.Request) {

}

func MessageRouter() *chi.Mux {
	router := chi.NewMux()

	router.Post("/", handleCreateMessage)

	return router
}
