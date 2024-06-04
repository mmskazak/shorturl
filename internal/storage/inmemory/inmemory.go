package inmemory

import (
	"errors"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"sync"
)

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
		return "", storageErrors.ErrNotFound
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
		return storageErrors.ErrKeyAlreadyExists
	}
	m.data[id] = targetURL
	return nil
}

func (m *InMemory) SaveBatch(items []storage.Incoming, baseHost string) ([]storage.Output, error) {
	result := make([]storage.Output, len(items))
	return result, nil
}
