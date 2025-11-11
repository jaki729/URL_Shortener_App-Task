/*
AI GENERATED

COPILOIT GENERATED FILE FOR TEST CASE GENERATION PURPOSES
TESTING COMPLETE FUNCTIONALITY OF THE CODE. DO NOT EDIT.

Package handler provides HTTP handlers for the URL shortener service.

This file contains unit tests for the handler package, covering the following scenarios:

- TestShortenURL_Success: Tests successful shortening of a valid URL, ensuring the response contains the expected short and long URLs.
- TestShortenURL_EmptyURL: Tests handling of an empty URL in the request, expecting a 400 Bad Request response.
- TestShortenURL_InvalidURL: Tests handling of an invalid URL format, expecting a 400 Bad Request response with an error message.
- TestShortenURL_InvalidJSON: Tests handling of invalid JSON in the request body, expecting a 400 Bad Request response.
- TestRedirectURL_Success: Tests successful redirection from a valid short code to the original long URL, expecting a 302 Found response with the correct Location header.
- TestRedirectURL_NotFound: Tests redirection with a non-existent short code, expecting a 404 Not Found response.
- TestShortenURL_Idempotency: Tests that shortening the same URL multiple times returns the same short URL, ensuring idempotency.

The tests use an in-memory storage implementation for isolation and do not require external dependencies.
*/
package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"URL_Shortener_Ruckus_Networks/internals/handler"
	"URL_Shortener_Ruckus_Networks/internals/service"
	"URL_Shortener_Ruckus_Networks/internals/storage"

	"github.com/gorilla/mux"
)

func setupHandler() *handler.Handler {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")
	return handler.NewHandler(svc)
}

func TestShortenURL_Success(t *testing.T) {
	h := setupHandler()
	reqBody := handler.ShortenRequest{
		URL: "https://www.example.com/test",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ShortenURL(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	var resp handler.ShortenResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.ShortURL == "" {
		t.Fatal("Expected short URL in response")
	}

	if resp.LongURL != reqBody.URL {
		t.Fatalf("Expected long URL %s, got %s", reqBody.URL, resp.LongURL)
	}
}

func TestShortenURL_EmptyURL(t *testing.T) {
	h := setupHandler()
	reqBody := handler.ShortenRequest{
		URL: "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ShortenURL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestShortenURL_InvalidURL(t *testing.T) {
	h := setupHandler()
	reqBody := handler.ShortenRequest{
		URL: "not-a-valid-url",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ShortenURL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
	var errResp handler.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Fatal("Expected error message in response")
	}
}

func TestShortenURL_InvalidJSON(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ShortenURL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestRedirectURL_Success(t *testing.T) {
	h := setupHandler()

	longURL := "https://www.example.com/test"
	reqBody := handler.ShortenRequest{URL: longURL}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ShortenURL(w, req)
	var shortenResp handler.ShortenResponse
	json.NewDecoder(w.Body).Decode(&shortenResp)

	// Extract short code from response
	shortCode := shortenResp.ShortURL[len("http://localhost:8080/"):]

	// Now test redirect
	req = httptest.NewRequest("GET", "/"+shortCode, nil)
	w = httptest.NewRecorder()

	// Use mux to parse URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/{shortCode}", h.RedirectURL)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("Expected status 302, got %d", w.Code)
	}

	location := w.Header().Get("Location")
	if location != longURL {
		t.Fatalf("Expected redirect to %s, got %s", longURL, location)
	}
}

func TestRedirectURL_NotFound(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	// Use mux to parse URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/{shortCode}", h.RedirectURL)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d", w.Code)
	}
}

func TestShortenURL_Idempotency(t *testing.T) {
	h := setupHandler()
	longURL := "https://www.example.com/idempotent-test"
	reqBody := handler.ShortenRequest{URL: longURL}
	body, _ := json.Marshal(reqBody)

	// First request
	req1 := httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	h.ShortenURL(w1, req1)
	var resp1 handler.ShortenResponse
	json.NewDecoder(w1.Body).Decode(&resp1)

	// Second request with same URL
	body, _ = json.Marshal(reqBody)
	req2 := httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	h.ShortenURL(w2, req2)
	var resp2 handler.ShortenResponse
	json.NewDecoder(w2.Body).Decode(&resp2)

	if resp1.ShortURL != resp2.ShortURL {
		t.Fatalf("Expected same short URL, got %s and %s", resp1.ShortURL, resp2.ShortURL)
	}
}
