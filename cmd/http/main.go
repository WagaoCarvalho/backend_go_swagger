package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WagaoCarvalho/backendgoswagger/cmd/http/routes"
)

func main() {
	port := "3000"
	fmt.Printf("Api running in port %s\n", port)
	r := routes.NewRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

}
