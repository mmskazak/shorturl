package mapstorage

import (
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/services/rwstorage"
	"strconv"
	"sync"
)

var ErrNotFound = errors.New("key not found")
var ErrKeyAlreadyExists = errors.New("key already exists")

type MapStorage struct {
	mu       *sync.Mutex
	Data     map[string]string
	FilePath string
}

func NewMapStorage(filePath string) *MapStorage {
	return &MapStorage{
		mu:       &sync.Mutex{},
		Data:     make(map[string]string),
		FilePath: filePath,
	}
}

func (m *MapStorage) GetShortURL(id string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	targetURL, ok := m.Data[id]
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
	if _, ok := m.Data[id]; ok {
		return ErrKeyAlreadyExists
	}
	m.Data[id] = targetURL

	if m.FilePath != "" {
		producer, err := rwstorage.NewProducer(m.FilePath)
		if err != nil {
			return fmt.Errorf("ошибка создания producer %w", err)
		}

		shData := rwstorage.ShortURLStruct{
			UUID:        strconv.Itoa(len(m.Data)),
			ShortURL:    id,
			OriginalURL: targetURL,
		}

		err = producer.WriteData(&shData)
		if err != nil {
			return fmt.Errorf("ошибка записи строки в файл %w", err)
		}
		producer.Close()
		logger.Logf.Infof("Добавлени которкая ссылка %v", shData)
	}
	return nil
}
