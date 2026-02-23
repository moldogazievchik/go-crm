package main

import (
	"log"
	"net/http"

	"github.com/moldogazievchik/go-crm/internal/httpapi"
)

func main() {
	handler := httpapi.Routes()

	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
