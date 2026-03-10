package main

import (
	"log"
	"net/http"
	"os"

	"github.com/moldogazievchik/go-crm/internal/httpapi"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, httpapi.RoutesWithConfig(httpapi.LoadConfig())))
}
