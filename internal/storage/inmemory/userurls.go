package inmemory

import (
	"context"
	"mmskazak/shorturl/internal/storage"
)

// GetUserURLs - получение всех URL для конкретного пользователя.
func (m *InMemory) GetUserURLs(ctx context.Context, userID string) ([]storage.URL, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var urls []storage.URL
	ids, ok := m.userIndex[userID]
	if !ok {
		return urls, nil // Пустой список, если пользователь не найден
	}

	for _, id := range ids {
		record, ok := m.data[id]
		if ok && !record.Deleted {
			urls = append(urls, storage.URL{
				ShortURL:    id,
				OriginalURL: record.OriginalURL,
			})
		}
	}

	return urls, nil
}
