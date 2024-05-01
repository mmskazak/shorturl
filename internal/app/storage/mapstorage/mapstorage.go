package mapstorage

import (
	"errors"
	"mmskazak/shorturl/internal/app/storage/storage"
	"sync"
)

var (
	instance *MapStorage
	once     sync.Once
)

type MapStorage struct {
	mu   sync.Mutex
	data map[string]string
}

func GetMapStorageInstance() *MapStorage {
	once.Do(func() {
		instance = &MapStorage{
			data: make(map[string]string),
		}
	})
	return instance
}

// SetMapStorageInstance устанавливает указанный экземпляр MapStorage.
func SetMapStorageInstance(ms *MapStorage) {
	instance = ms
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
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[id]; ok {
		return errors.New("key already exists")
	}
	m.data[id] = targetURL
	return nil
}
