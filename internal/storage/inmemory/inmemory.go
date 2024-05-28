package inmemory

import (
	"errors"
	"sync"
)

var ErrKeyAlreadyExists = errors.New("key already exists")
var ErrNotFound = errors.New("key not found")

type InMemory struct {
	mu   *sync.Mutex
	data map[string]string
}

// NumberOfEntries - количество записей.
func (m *InMemory) NumberOfEntries() int {
	return len(m.data)
}

func NewInMemory() (*InMemory, error) {
	return &InMemory{
		mu:   &sync.Mutex{},
		data: make(map[string]string),
	}, nil
}

func (m *InMemory) GetShortURL(id string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	targetURL, ok := m.data[id]
	if !ok {
		return "", ErrNotFound
	}
	return targetURL, nil
}

func (m *InMemory) SetShortURL(id string, targetURL string) error {
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
