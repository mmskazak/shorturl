package contracts

import (
	"context"
	"mmskazak/shorturl/internal/models"
)

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
		items []models.Incoming,
		baseHost string,
		userID string,
		generator IGenIDForURL,
	) ([]models.Output, error)
	GetUserURLs(ctx context.Context, userID string, baseHost string) ([]models.URL, error)
	DeleteURLs(ctx context.Context, urlIDs []string) error
}
