package inmemory

import (
	"context"
	"fmt"
	"net/url"

	"mmskazak/shorturl/internal/storage"
)

// GetUserURLs - получение всех URL для конкретного пользователя.
//
// Функция возвращает список сокращенных и оригинальных URL-адресов для заданного пользователя.
// Если пользователь не найден, возвращается пустой список.
// Если возникают ошибки при парсинге базового хоста или относительного пути, возвращается ошибка.
//
// Параметры:
// - ctx: контекст выполнения запроса.
// - userID: идентификатор пользователя, для которого нужно получить URL-адреса.
// - baseHost: базовый хост для создания полного сокращенного URL.
//
// Возвращаемые значения:
// - []storage.URL: список структур URL, содержащих сокращенные и оригинальные URL-адреса.
// - error: ошибка, если она произошла в процессе выполнения запроса.
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
