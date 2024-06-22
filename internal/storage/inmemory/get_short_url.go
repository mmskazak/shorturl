package inmemory

import (
	"context"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

// GetShortURL - получение оригинального URL по короткому идентификатору.
func (m *InMemory) GetShortURL(_ context.Context, id string) (string, error) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	record, ok := m.Data[id]
	if !ok {
		return "", storageErrors.ErrNotFound
	}
	if record.Deleted {
		return "", storageErrors.ErrDeletedShortURL
	}

	return record.OriginalURL, nil
}
