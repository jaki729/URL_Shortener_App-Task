/*
AI GENERATED

COPILOIT GENERATED FILE FOR TEST CASE GENERATION PURPOSES
TESTING COMPLETE FUNCTIONALITY OF THE CODE. DO NOT EDIT.

This file contains unit tests and benchmarks for the core URL shortening service in the URL shortener application.

The tests cover the following aspects:
- Successful URL shortening and retrieval, ensuring that valid URLs are correctly shortened and can be expanded back to their original form.
- Idempotency of the shortening operation, verifying that shortening the same URL multiple times yields the same short code and short URL.
- Validation of input URLs, including handling of empty strings, unsupported schemes, and missing schemes, as well as acceptance of valid HTTP and HTTPS URLs.
- Retrieval of original URLs from short codes, including correct handling of non-existent short codes and expected error responses.
- Deterministic and unique short code generation, ensuring that the same URL always produces the same code, and different URLs produce different codes.
- Performance benchmarks for both the shortening and retrieval operations, providing insights into the efficiency of the service under repeated use.

These tests ensure the correctness, reliability, and performance of the URLService, which is responsible for the main business logic of the URL shortener application.
*/

package test

import (
	"strings"
	"testing"

	"URL_Shortener_Ruckus_Networks/internals/service"
	"URL_Shortener_Ruckus_Networks/internals/storage"
)

func TestService_ShortURL_Success(t *testing.T) {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")

	longURL := "https://www.example.com/very/long/url/path"
	shortURL, shortCode, err := svc.ShortenURL(longURL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if shortURL == "" {
		t.Fatal("Expected short URL, got empty string")
	}

	if shortCode == "" {
		t.Fatal("Expected short code, got empty string")
	}

	if !strings.HasPrefix(shortURL, "http://localhost:8080/") {
		t.Fatalf("Expected short URL to start with base URL, got %s", shortURL)
	}
}

func TestService_ShortURL_Idempotency(t *testing.T) {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")

	longURL := "https://www.example.com/test"

	// First call
	shortURL1, shortCode1, err1 := svc.ShortenURL(longURL)
	if err1 != nil {
		t.Fatalf("First call failed: %v", err1)
	}

	// Second call with same URL
	shortURL2, shortCode2, err2 := svc.ShortenURL(longURL)
	if err2 != nil {
		t.Fatalf("Second call failed: %v", err2)
	}

	// Should return the same short URL
	if shortURL1 != shortURL2 {
		t.Fatalf("Expected same short URL, got %s and %s", shortURL1, shortURL2)
	}

	if shortCode1 != shortCode2 {
		t.Fatalf("Expected same short code, got %s and %s", shortCode1, shortCode2)
	}
}

func TestService_ShortURL_InvalidURL(t *testing.T) {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")

	testCases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Empty URL", "", true},
		{"Invalid scheme", "ftp://example.com", true},
		{"No scheme", "example.com", true},
		{"Valid HTTP", "http://example.com", false},
		{"Valid HTTPS", "https://example.com", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := svc.ShortenURL(tc.url)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestGetLongURL_Success(t *testing.T) {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")

	longURL := "https://www.example.com/test"
	_, shortCode, err := svc.ShortenURL(longURL)
	if err != nil {
		t.Fatalf("Failed to shorten URL: %v", err)
	}

	retrievedURL, err := svc.GetLongURL(shortCode)
	if err != nil {
		t.Fatalf("Failed to get long URL: %v", err)
	}

	if retrievedURL != longURL {
		t.Fatalf("Expected %s, got %s", longURL, retrievedURL)
	}
}

func TestGetLongURL_NotFound(t *testing.T) {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")

	_, err := svc.GetLongURL("nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent short code")
	}

	if err != storage.ErrNotFound {
		t.Fatalf("Expected ErrNotFound, got %v", err)
	}
}

func TestGenerateShortCode_Deterministic(t *testing.T) {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")

	url := "https://www.example.com/test"

	code1 := svc.GenerateShortCode(url)
	code2 := svc.GenerateShortCode(url)

	if code1 != code2 {
		t.Fatalf("Expected deterministic short codes, got %s and %s", code1, code2)
	}
}

func TestGenerateShortCode_Unique(t *testing.T) {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")

	url1 := "https://www.example.com/test1"
	url2 := "https://www.example.com/test2"

	code1 := svc.GenerateShortCode(url1)
	code2 := svc.GenerateShortCode(url2)

	if code1 == code2 {
		t.Fatalf("Expected different short codes for different URLs")
	}
}

func BenchmarkShortenURL(b *testing.B) {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")
	longURL := "https://www.example.com/benchmark"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := svc.ShortenURL(longURL)
		if err != nil {
			b.Fatalf("ShortenURL failed: %v", err)
		}
	}
}

func BenchmarkGetLongURL(b *testing.B) {
	store := storage.NewMemoryStorage()
	svc := service.NewURLService(store, "http://localhost:8080")
	longURL := "https://www.example.com/benchmark"
	_, shortCode, err := svc.ShortenURL(longURL)
	if err != nil {
		b.Fatalf("ShortenURL failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.GetLongURL(shortCode)
		if err != nil {
			b.Fatalf("GetLongURL failed: %v", err)
		}
	}
}
