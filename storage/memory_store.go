package storage

import (
	"sync"
	"time"
)

type URLData struct {
	LongURL   string
	ExpiresAt time.Time
}

type MemoryStore struct {
	store map[string]URLData
	mu    sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]URLData),
	}
}

func (m *MemoryStore) Save(shortURL, longURL string, expirationSeconds int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	expirationTime := time.Now().Add(time.Duration(expirationSeconds) * time.Second)

	m.store[shortURL] = URLData{
		LongURL:   longURL,
		ExpiresAt: expirationTime,
	}
}

func (m *MemoryStore) Get(shortURL string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, exists := m.store[shortURL]
	if !exists {
		return "", false
	}

	if time.Now().After(data.ExpiresAt) {
		delete(m.store, shortURL)
		return "", false
	}

	return data.LongURL, true
}
