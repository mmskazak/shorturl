package inmemory

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/storage"
	"net/url"
)

// GetUserURLs - получение всех URL для конкретного пользователя.
func (m *InMemory) GetUserURLs(ctx context.Context, userID string, baseHost string) ([]storage.URL, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var urls []storage.URL
	ids, ok := m.userIndex[userID]
	if !ok {
		return urls, nil // Пустой список, если пользователь не найден
	}

	// Парсим базовый хост
	baseURL, err := url.Parse(baseHost)
	if err != nil {
		return nil, fmt.Errorf("error parsing baseHost: %w", err)
	}

	for _, id := range ids {
		record, ok := m.data[id]
		if ok && !record.Deleted {
			// Парсим относительный путь (id)
			relativeURL, err := url.Parse(id)
			if err != nil {
				return nil, fmt.Errorf("error parsing relative URL: %w", err)
			}

			// Объединяем baseURL и relativeURL
			fullShortURL := baseURL.ResolveReference(relativeURL).String()

			urls = append(urls, storage.URL{
				ShortURL:    fullShortURL,
				OriginalURL: record.OriginalURL,
			})
		}
	}

	return urls, nil
}
