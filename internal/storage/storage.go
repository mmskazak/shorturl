package storage

import (
	"context"
	"fmt"
	"net/url"
)

// Incoming представляет данные, полученные при создании или обновлении короткого URL.
// CorrelationID - идентификатор запроса, связанный с этим URL.
// OriginalURL - оригинальный (длинный) URL, который нужно сократить.
type Incoming struct {
	CorrelationID string `json:"correlation_id"` // строковый идентификатор
	OriginalURL   string `json:"original_url"`   // оригинальный URL
}

// Output представляет результат операции создания или обновления короткого URL.
// CorrelationID - идентификатор запроса, связанный с этим URL.
// ShortURL - сгенерированный короткий URL.
type Output struct {
	CorrelationID string `json:"correlation_id"` // строковый идентификатор
	ShortURL      string `json:"short_url"`      // короткий URL
}

// URL представляет собой короткий URL и его оригинальный (длинный) URL.
// ShortURL - сгенерированный короткий URL.
// OriginalURL - оригинальный URL.
type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// IGenIDForURL представляет интерфейс для генерации идентификаторов для коротких URL.
type IGenIDForURL interface {
	Generate() (string, error) // Метод для генерации нового идентификатора.
}

// Storage представляет интерфейс для взаимодействия со хранилищем коротких URL-адресов.
// Close - закрывает соединение с хранилищем.
// GetShortURL - получает оригинальный URL по короткому URL.
// SetShortURL - устанавливает связь между коротким URL и оригинальным URL.
// SaveBatch - сохраняет пакет новых коротких URL-адресов.
// GetUserURLs - получает все короткие URL-адреса, связанные с пользователем.
// DeleteURLs - устанавливает флаг удаления для указанного списка URL.
type Storage interface {
	Close() error
	GetShortURL(ctx context.Context, id string) (string, error)
	SetShortURL(ctx context.Context, idShortPath string, targetURL string, userID string, deleted bool) error
	SaveBatch(
		ctx context.Context,
		items []Incoming,
		baseHost string,
		userID string,
		generator IGenIDForURL,
	) ([]Output, error)
	GetUserURLs(ctx context.Context, userID string, baseHost string) ([]URL, error)
	DeleteURLs(ctx context.Context, urlIDs []string) error
}

// GetFullShortURL создает полный короткий URL, используя базовый хост и относительный путь (короткий URL).
// BaseHost - базовый хост для создания полного URL.
// CorrelationID - относительный путь, который будет добавлен к базовому хосту.
// Возвращает полный короткий URL и ошибку, если она произошла.
func GetFullShortURL(baseHost, correlationID string) (string, error) {
	u, err := url.Parse(baseHost)
	if err != nil {
		return "", fmt.Errorf("error parsing baseHost: %w", err)
	}
	// ResolveReference корректно объединяет базовый URL и относительный путь
	u = u.ResolveReference(&url.URL{Path: correlationID})
	return u.String(), nil
}
