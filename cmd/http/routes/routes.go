package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backendgoswagger/internal/utils/handlers"
)

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /", handlers.GetHome)
	return r
}
