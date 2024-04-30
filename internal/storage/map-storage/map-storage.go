package map_storage

import (
	"errors"
	"mmskazak/shorturl/internal/storage/storage"
	"sync"
)

var (
	instance *MapStorage
	once     sync.Once
)

type Storage interface {
	GetShortURL(id string) string            //получение короткого URL
	SetShortURL(id string, targetURL string) //установить короткий URL
}

type MapStorage struct {
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
	targetURL, ok := m.data[id]
	if ok != true {
		return "", storage.ErrNotFound
	}
	return targetURL, nil
}

func (m *MapStorage) SetShortURL(id string, targetURL string) error {
	if _, ok := m.data[id]; ok == true {
		return errors.New("key already exists")
	}
	m.data[id] = targetURL
	return nil
}
