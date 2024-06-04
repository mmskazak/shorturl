package inmemory

import (
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"sync"
)

type InMemory struct {
	mu           *sync.Mutex
	data         map[string]string
	indexForData map[string]string
}

// NumberOfEntries - количество записей.
func (m *InMemory) NumberOfEntries() int {
	return len(m.data)
}

func NewInMemory() (*InMemory, error) {
	return &InMemory{
		mu:           &sync.Mutex{},
		data:         make(map[string]string),
		indexForData: make(map[string]string),
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
	if _, ok := m.indexForData[targetURL]; ok {
		return storageErrors.ErrOriginalURLAlreadyExists
	}
	m.data[id] = targetURL
	m.indexForData[targetURL] = id
	return nil
}

func (m *InMemory) SaveBatch(items []storage.Incoming, baseHost string) ([]storage.Output, error) {
	dontChangedData := m.data

	outputs := make([]storage.Output, 0, len(items))
	for _, v := range items {
		err := m.SetShortURL(v.CorrelationID, v.OriginalURL)
		if err != nil {
			m.data = dontChangedData
			return nil, fmt.Errorf("save batch error: %w", err)
		}

		fullShortURL, err := storage.GetFullShortURL(baseHost, v.CorrelationID)
		if err != nil {
			m.data = dontChangedData
			return nil, fmt.Errorf("error getFullShortURL from two parts %w", err)
		}

		outputs = append(outputs, storage.Output{
			CorrelationID: v.CorrelationID,
			ShortURL:      fullShortURL,
		})
	}

	return outputs, nil
}
