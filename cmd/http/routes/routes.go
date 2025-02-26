package routes

import (
	"github.com/WagaoCarvalho/backendgoswagger/internal/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	// Cria um novo roteador com gorilla/mux
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", handlers.GetHome).Methods("GET")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/user", handlers.PostUser).Methods("POST")
	r.HandleFunc("/user/{id}", handlers.PutUser).Methods("PUT")
	r.HandleFunc("/user/{id}", handlers.DeleteUser).Methods("DELETE")

	return r
}
