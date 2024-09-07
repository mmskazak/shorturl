package storage

import (
	"fmt"
	"net/url"
)

// GetFullShortURL создает полный короткий URL, используя базовый хост и относительный путь (короткий URL).
// BaseHost - базовый хост для создания полного URL.
// CorrelationID - относительный путь, который будет добавлен к базовому хосту.
// Возвращает полный короткий URL и ошибку, если она произошла.
func GetFullShortURL(baseHost, shortURL string) (string, error) {
	u, err := url.Parse(baseHost)
	if err != nil {
		return "", fmt.Errorf("error parsing baseHost: %w", err)
	}
	// ResolveReference корректно объединяет базовый URL и относительный путь
	u = u.ResolveReference(&url.URL{Path: shortURL})
	return u.String(), nil
}
