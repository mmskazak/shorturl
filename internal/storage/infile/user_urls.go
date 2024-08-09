package infile

import (
	"context"
	"mmskazak/shorturl/internal/storage"
)

// GetUserURLs - получение всех URL для конкретного пользователя.
func (m *InFile) GetUserURLs(ctx context.Context, userID string, baseHost string) ([]storage.URL, error) {
	urls, err := m.InMe.GetUserURLs(ctx, userID, baseHost)
	if err != nil {
		return nil, err //nolint:wrapcheck //пробрасываем дальше ошибку
	}

	return urls, nil
}
