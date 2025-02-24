package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backendgoswagger/internal/handlers"
)

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /", handlers.GetHome)
	r.HandleFunc("GET /users", handlers.GetUsers)
	r.HandleFunc("POST /user", handlers.PostUser)
	return r
}
