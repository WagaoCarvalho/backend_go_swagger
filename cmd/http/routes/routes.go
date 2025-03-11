package routes

import (
	"github.com/WagaoCarvalho/backendgoswagger/internal/handlers"
	"github.com/WagaoCarvalho/backendgoswagger/middlewares"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(middlewares.Logging)
	r.Use(middlewares.RecoverPanic)
	r.Use(middlewares.RateLimiter)
	r.Use(middlewares.RateLimiter)
	r.Use(middlewares.CORS)

	r.HandleFunc("/", handlers.GetHome).Methods("GET")
	r.HandleFunc("/user", handlers.PostUser).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	protectedRoutes := r.PathPrefix("/").Subrouter()
	protectedRoutes.Use(middlewares.IsAuthByBearerToken)

	protectedRoutes.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")
	protectedRoutes.HandleFunc("/user/{id}", handlers.PutUser).Methods("PUT")
	protectedRoutes.HandleFunc("/user/{id}", handlers.DeleteUser).Methods("DELETE")
	protectedRoutes.HandleFunc("/users", handlers.GetUsers).Methods("GET")

	protectedRoutes.HandleFunc("/wallets", handlers.GetWallets).Methods("GET")
	protectedRoutes.HandleFunc("/wallet", handlers.GetWallet).Methods("GET")
	protectedRoutes.HandleFunc("/wallet", handlers.PutWallet).Methods("PUT")

	protectedRoutes.HandleFunc("/transaction", handlers.PostTransaction).Methods("POST")
	protectedRoutes.HandleFunc("/transactions", handlers.GetTransactions).Methods("GET")

	return r
}
