package inmemory

import (
	"context"
	"errors"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"sync"

	"go.uber.org/zap"
)

type InMemory struct {
	mu           *sync.Mutex
	data         map[string]string
	indexForData map[string]string
	zapLog       *zap.SugaredLogger
}

// NumberOfEntries - количество записей.
func (m *InMemory) NumberOfEntries() int {
	return len(m.data)
}

func NewInMemory(zapLog *zap.SugaredLogger) (*InMemory, error) {
	return &InMemory{
		mu:           &sync.Mutex{},
		data:         make(map[string]string),
		indexForData: make(map[string]string),
		zapLog:       zapLog,
	}, nil
}

func (m *InMemory) GetShortURL(_ context.Context, id string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	targetURL, ok := m.data[id]
	if !ok {
		return "", storageErrors.ErrNotFound
	}
	return targetURL, nil
}

func (m *InMemory) SetShortURL(_ context.Context, id string, targetURL string) error {
	if id == "" {
		return errors.New("id is empty")
	}
	if targetURL == "" {
		return errors.New("URL is empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.indexForData[targetURL]; ok {
		return storageErrors.ConflictError{
			ShortURL: id,
			Err:      storageErrors.ErrOriginalURLAlreadyExists,
		}
	}
	if _, ok := m.data[id]; ok {
		return storageErrors.ErrKeyAlreadyExists
	}
	m.data[id] = targetURL
	m.indexForData[targetURL] = id
	return nil
}

func (m *InMemory) Close() error {
	// На данный момент закрывать нечего, но метод оставлен для возможных будущих изменений
	m.zapLog.Debugln("InMemory storage closed (nothing to close currently)")
	return nil
}
