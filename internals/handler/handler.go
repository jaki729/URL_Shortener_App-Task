package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"URL_Shortener_Ruckus_Networks/internals/service"
	"URL_Shortener_Ruckus_Networks/internals/storage"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.URLService
}

// creating new handler instance
func NewHandler(service *service.URLService) *Handler {
	return &Handler{
		service: service,
	}
}

// ShortenRequest handler
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse handler
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

// ErrorResponse handle
type ErrorResponse struct {
	Error string `json:"error"`
}

// ShortenURL API - POST /api/shorten
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("handler: ShortenURL - invalid request body: %v", err)
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("handler: ShortenURL - incoming URL=%s", req.URL)

	if req.URL == "" {
		log.Printf("handler: ShortenURL - empty URL")
		h.sendError(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortURL, _, err := h.service.ShortenURL(req.URL)
	if err != nil {
		if err == service.ErrInvalidURL {
			log.Printf("handler: ShortenURL - invalid URL format: %s", req.URL)
			h.sendError(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		log.Printf("handler: ShortenURL - internal error: %v", err)
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := ShortenResponse{
		ShortURL: shortURL,
		LongURL:  req.URL,
	}

	log.Printf("handler: ShortenURL - created short_url=%s for long_url=%s", shortURL, req.URL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RedirectURL API - GET /{shortCode}
func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	log.Printf("handler: RedirectURL - shortCode=%s method=%s", shortCode, r.Method)

	if shortCode == "" {
		log.Printf("handler: RedirectURL - missing short code")
		h.sendError(w, "Short code is required", http.StatusBadRequest)
		return
	}

	longURL, err := h.service.GetLongURL(shortCode)
	if err != nil {
		if err == storage.ErrNotFound {
			log.Printf("handler: RedirectURL - not found shortCode=%s", shortCode)
			h.sendError(w, "Short URL not found", http.StatusNotFound)
			return
		}
		log.Printf("handler: RedirectURL - internal error retrieving shortCode=%s: %v", shortCode, err)
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("handler: RedirectURL - redirecting shortCode=%s -> %s", shortCode, longURL)
	http.Redirect(w, r, longURL, http.StatusFound)
}

// sendError
func (h *Handler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
