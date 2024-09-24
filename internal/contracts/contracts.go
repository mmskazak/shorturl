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

// IInternalStats - внутренняя статистика по приложению
type IInternalStats interface {
	InternalStats(ctx context.Context) (models.Stats, error)
}
