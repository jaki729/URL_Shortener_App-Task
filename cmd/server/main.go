package main

import (
	"log"
	"net/http"
	"os"

	"URL_Shortener_Ruckus_Networks/internals/handler"
	"URL_Shortener_Ruckus_Networks/internals/service"
	"URL_Shortener_Ruckus_Networks/internals/storage"

	"github.com/gorilla/mux"
)

func main() {
	// Env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// storage
	store := storage.NewMemoryStorage()

	// service
	svc := service.NewURLService(store, baseURL)

	// handler
	h := handler.NewHandler(svc)

	// Routers
	r := mux.NewRouter()
	r.HandleFunc("/api/shorten", h.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortCode}", h.RedirectURL).Methods("GET", "HEAD")

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Printf("Base URL: %s", baseURL)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
