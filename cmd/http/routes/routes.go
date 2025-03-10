package routes

import (
	"github.com/WagaoCarvalho/backendgoswagger/internal/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", handlers.GetHome).Methods("GET")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/user", handlers.PostUser).Methods("POST")
	r.HandleFunc("/user/{id}", handlers.PutUser).Methods("PUT")
	r.HandleFunc("/user/{id}", handlers.DeleteUser).Methods("DELETE")

	r.HandleFunc("/wallets", handlers.GetWallets).Methods("GET")
	r.HandleFunc("/wallet", handlers.GetWallet).Methods("GET")
	r.HandleFunc("/wallet", handlers.PutWallet).Methods("PUT")

	r.HandleFunc("/transaction", handlers.PostTransaction).Methods("POST")
	r.HandleFunc("/transactions", handlers.GetTransactions).Methods("GET")

	r.HandleFunc("/login", handlers.Login).Methods("POST")
	return r
}
