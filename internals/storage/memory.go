package storage

import (
	"log"
	"sync"
)

// MemoryStorage implements Storage using in-memory maps
type MemoryStorage struct {
	shortToLong map[string]string
	longToShort map[string]string
	mu          sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortToLong: make(map[string]string),
		longToShort: make(map[string]string),
	}
}

// map of shortCode to longURL
func (m *MemoryStorage) Save(shortCode, longURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("storage: Save - shortCode=%s longURL=%s", shortCode, longURL)
	m.shortToLong[shortCode] = longURL
	m.longToShort[longURL] = shortCode

	return nil
}

// retrieve longURL by shortCode
func (m *MemoryStorage) GetLongURL(shortCode string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	longURL, exists := m.shortToLong[shortCode]
	if !exists {
		log.Printf("storage: GetLongURL - not found shortCode=%s", shortCode)
		return "", ErrNotFound
	}

	log.Printf("storage: GetLongURL - found shortCode=%s longURL=%s", shortCode, longURL)
	return longURL, nil
}

// retrieve shortCode by longURL
func (m *MemoryStorage) GetShortCode(longURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	shortCode, exists := m.longToShort[longURL]
	if !exists {
		log.Printf("storage: GetShortCode - not found longURL=%s", longURL)
		return "", ErrNotFound
	}

	log.Printf("storage: GetShortCode - found longURL=%s shortCode=%s", longURL, shortCode)
	return shortCode, nil
}

// check if shortCode exists
func (m *MemoryStorage) Exists(shortCode string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.shortToLong[shortCode]
	log.Printf("storage: Exists - shortCode=%s exists=%v", shortCode, exists)
	return exists
}
