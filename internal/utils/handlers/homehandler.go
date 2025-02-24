package handlers

import (
	"net/http"

	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	utils.ToJson(w, struct {
		Message string `json:"message"`
	}{
		Message: "Go RESTful Api x",
	})
}

//mux.HandleFunc("/", noteHandler.NotesList)

//mux.Handle("GET /note/{id}", handlers.HandlerWithError(noteHandler.NoteView))
