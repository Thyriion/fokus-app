package main

import (
	"fokus-app/internal/web"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Direkt Handler registrieren
	r.Get("/", web.HandleIndex)

	port := "8080"
	log.Printf("Starting server on port %s...\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
