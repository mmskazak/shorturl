package mapstorage

import (
	"errors"
	"mmskazak/shorturl/internal/app/storage/storage"
	"sync"
)

type MapStorage struct {
	mu   sync.Mutex
	data map[string]string
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		data: make(map[string]string),
	}
}

func (m *MapStorage) GetShortURL(id string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	targetURL, ok := m.data[id]
	if !ok {
		return "", storage.ErrNotFound
	}
	return targetURL, nil
}

func (m *MapStorage) SetShortURL(id string, targetURL string) error {
	if id == "" {
		return errors.New("id is empty")
	}
	if targetURL == "" {
		return errors.New("URL is empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[id]; ok {
		return errors.New("key already exists")
	}
	m.data[id] = targetURL
	return nil
}
