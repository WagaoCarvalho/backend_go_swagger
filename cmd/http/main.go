package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WagaoCarvalho/backendgoswagger/cmd/http/routes"
	"github.com/WagaoCarvalho/backendgoswagger/config"
	"github.com/WagaoCarvalho/backendgoswagger/internal/repositories"
)

func main() {
	configs := config.LoadConfig()
	port := configs.ServerPort
	if port == "" {
		port = "5000"
	}

	repositories.TestConnection()
	fmt.Printf("API running on port %s\n", port)

	r := routes.NewRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
