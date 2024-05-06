package mapstorage

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("key not found")
var ErrKeyAlreadyExists = errors.New("key already exists")

type MapStorage struct {
	mu   *sync.Mutex
	data map[string]string
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		mu:   &sync.Mutex{},
		data: make(map[string]string),
	}
}

func (m *MapStorage) GetShortURL(id string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	targetURL, ok := m.data[id]
	if !ok {
		return "", ErrNotFound
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
		return ErrKeyAlreadyExists
	}
	m.data[id] = targetURL
	return nil
}
