package inmemory

import (
	"context"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

// GetShortURL - получение оригинального URL по короткому идентификатору.
func (m *InMemory) GetShortURL(_ context.Context, id string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	record, ok := m.data[id]
	if !ok || record.Deleted {
		return "", storageErrors.ErrNotFound
	}
	return record.OriginalURL, nil
}
