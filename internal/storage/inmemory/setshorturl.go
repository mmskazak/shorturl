package inmemory

import (
	"context"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

// SetShortURL error:
// different error
// ErrKeyAlreadyExists
// ConflictError (ErrOriginalURLAlreadyExists)
func (m *InMemory) SetShortURL(_ context.Context, id string, originalURL string, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Проверка на существование URL.
	for _, record := range m.data {
		if record.OriginalURL == originalURL && !record.Deleted {
			return storageErrors.ConflictError{
				ShortURL: id,
				Err:      storageErrors.ErrOriginalURLAlreadyExists,
			}
		}
	}

	// Проверка на существование id.
	if _, ok := m.data[id]; ok {
		return storageErrors.ErrKeyAlreadyExists
	}

	// Добавление записи.
	m.data[id] = URLRecord{
		OriginalURL: originalURL,
		UserID:      userID,
		Deleted:     false,
	}
	m.userIndex[userID] = append(m.userIndex[userID], id)

	return nil
}
