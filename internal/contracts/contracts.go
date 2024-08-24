package contracts

import (
	"context"
	"mmskazak/shorturl/internal/dtos"

	"mmskazak/shorturl/internal/models"
)

//go:generate mockgen -source=contracts.go -destination=mocks/contracts.go -package=mocks

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
		items []models.Incoming,
		baseHost string,
		userID string,
		generator IGenIDForURL,
	) ([]models.Output, error)
	GetUserURLs(ctx context.Context, userID string, baseHost string) ([]models.URL, error)
	DeleteURLs(ctx context.Context, urlIDs []string) error
}

// Pinger определяет интерфейс для проверки состояния хранилища.
type Pinger interface {
	Ping(ctx context.Context) error
}

// ISaveBatch сохранение URLs батчем.
type ISaveBatch interface {
	SaveBatch(
		ctx context.Context,
		items []models.Incoming,
		baseHost string,
		userID string,
		generator IGenIDForURL,
	) ([]models.Output, error)
}

// ISetShortURL устанавливает связь между коротким URL и оригинальным URL, сохраняет в хранилище.
type ISetShortURL interface {
	SetShortURL(ctx context.Context, idShortPath string, targetURL string, userID string, deleted bool) error
}

// IDeleteUserURLs устанавливает флаг удаления для множества записей в хранилище.
type IDeleteUserURLs interface {
	DeleteURLs(ctx context.Context, urlIDs []string) error
}

// IGetUserURLs возвращает все URL-адреса, связанные с указанным пользователем.
type IGetUserURLs interface {
	GetUserURLs(ctx context.Context, userID string, baseHost string) ([]models.URL, error)
}

// IGetShortURL - получение оригинального URL по короткому идентификатору.
type IGetShortURL interface {
	GetShortURL(ctx context.Context, idShortPath string) (string, error)
}

// IShortURLService описывает контракт для сервиса создания и управления короткими URL.
type IShortURLService interface {
	// GenerateShortURL создает короткий URL, используя данные из DTO и генератор ID.
	GenerateShortURL(
		ctx context.Context,
		dto dtos.DTOShortURL,
		generator IGenIDForURL,
		data ISetShortURL,
	) (string, error)
}
