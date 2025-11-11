package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"URL_Shortener_Ruckus_Networks/internals/storage"
)

var (
	ErrInvalidURL = errors.New("invalid URL format")
)

// business logic struct for URL shortening service
type URLService struct {
	storage storage.Storage
	baseURL string
}

// creates a new URL service
func NewURLService(storage storage.Storage, baseURL string) *URLService {
	return &URLService{
		storage: storage,
		baseURL: baseURL,
	}
}

// idempotent receiver method - same long URL always returns same short URL
func (s *URLService) ShortenURL(longURL string) (string, string, error) {

	if err := s.validateURL(longURL); err != nil {
		log.Printf("service: ShortenURL - invalid URL=%s", longURL)
		return "", "", err
	}

	// idempotency check - return existing short code if present
	if shortCode, err := s.storage.GetShortCode(longURL); err == nil {
		shortURL := fmt.Sprintf("%s/%s", s.baseURL, shortCode)
		log.Printf("service: ShortenURL - existing mapping found longURL=%s shortCode=%s", longURL, shortCode)
		return shortURL, shortCode, nil
	}

	shortCode := s.GenerateShortCode(longURL)
	log.Printf("service: ShortenURL - generated shortCode=%s for longURL=%s", shortCode, longURL)

	if err := s.storage.Save(shortCode, longURL); err != nil {
		log.Printf("service: ShortenURL - failed to save mapping shortCode=%s longURL=%s err=%v", shortCode, longURL, err)
		return "", "", err
	}

	shortURL := fmt.Sprintf("%s/%s", s.baseURL, shortCode)
	log.Printf("service: ShortenURL - saved mapping shortCode=%s shortURL=%s", shortCode, shortURL)
	return shortURL, shortCode, nil
}

// get long URL by short code
func (s *URLService) GetLongURL(shortCode string) (string, error) {
	longURL, err := s.storage.GetLongURL(shortCode)
	if err != nil {
		log.Printf("service: GetLongURL - not found shortCode=%s err=%v", shortCode, err)
		return "", err
	}
	log.Printf("service: GetLongURL - found shortCode=%s longURL=%s", shortCode, longURL)
	return longURL, nil
}

// URL validation
func (s *URLService) validateURL(urlStr string) error {
	if urlStr == "" {
		return ErrInvalidURL
	}

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return ErrInvalidURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidURL
	}

	return nil
}

// code generation using SHA-256 and base64 encoding
func (s *URLService) GenerateShortCode(longURL string) string {
	log.Printf("service: GenerateShortCode - generating for URL=%s", longURL)
	hash := sha256.Sum256([]byte(longURL))
	encoded := base64.URLEncoding.EncodeToString(hash[:])

	shortCode := strings.TrimRight(encoded, "=")[:8]

	return shortCode
}
