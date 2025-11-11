package storage

import "errors"

var (
	ErrNotFound      = errors.New("short code not found")
	ErrAlreadyExists = errors.New("mapping already exists")
)

// interface for URL storage
type Storage interface {

	// map of shortCode to longURL
	Save(shortCode, longURL string) error

	// retrieve longURL by shortCode
	GetLongURL(shortCode string) (string, error)

	// retrieve shortCode by longURL
	GetShortCode(longURL string) (string, error)

	// check if shortCode exists
	Exists(shortCode string) bool
}