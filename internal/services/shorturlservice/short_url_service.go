package shorturlservice

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

//go:generate mockgen -source=short_url_service.go -destination=mocks/mock_short_url_service.go -package=mocks

// ErrServiceGenerateID указывает на ошибку, связанную с генерацией уникального идентификатора для короткого URL.
// Эта ошибка может возникнуть, если система не может создать новый уникальный идентификатор.
var ErrServiceGenerateID = errors.New("generateID failed")

// ErrConflict указывает на ошибку, когда оригинальный URL уже существует в системе.
// Эта ошибка возникает, если попытка сохранить короткий URL не удалась из-за конфликта с уже существующим URL.
var ErrConflict = errors.New("error original url already exists")

// IGenIDForURL описывает интерфейс для генерации идентификаторов коротких URL.
type IGenIDForURL interface {
	Generate() (string, error)
}

// DTOShortURL представляет данные, необходимые для создания короткого URL.
type DTOShortURL struct {
	OriginalURL string // Оригинальный URL
	UserID      string // Идентификатор пользователя
	BaseHost    string // Базовый хост для формирования короткого URL
	Deleted     bool   // Флаг, указывающий, удален ли URL
}

// ShortURLService предоставляет услуги по созданию и управлению короткими URL.
type ShortURLService struct {
	maxIteration int // Максимальное количество попыток генерации короткого URL
}

// ISetShortURL устанавливает связь между коротким URL и оригинальным URL, сохраняет в хранилище.
type ISetShortURL interface {
	SetShortURL(ctx context.Context, idShortPath string, targetURL string, userID string, deleted bool) error
}

// GenerateShortURL создает короткий URL, используя данные из DTO и генератор ID.
func (s *ShortURLService) GenerateShortURL(
	ctx context.Context,
	dto DTOShortURL,
	generator IGenIDForURL,
	data ISetShortURL,
) (string, error) {
	var err error
	// Разбираем базовый URL
	base, err := url.Parse(dto.BaseHost)
	if err != nil {
		return "", fmt.Errorf("ошибка при разборе базового URL: %w", err)
	}

	var id string
	for range s.maxIteration {
		// Генерируем новый идентификатор
		id, err = generator.Generate()
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrServiceGenerateID, err)
		}

		// Пытаемся сохранить короткий URL в хранилище
		err = data.SetShortURL(ctx, id, dto.OriginalURL, dto.UserID, dto.Deleted)
		if err != nil {
			// Проверяем, возникла ли ошибка конфликта
			conflictError, ok := isConflictError(err)
			if ok {
				// Парсим базовый хост
				baseURL, err := url.Parse(dto.BaseHost)
				if err != nil {
					return "", fmt.Errorf("не получилось распарсить dto base_host %w", err)
				}

				// Парсинг сокращенного URL
				shortURL, err := url.Parse(conflictError.ShortURL)
				if err != nil {
					return "", fmt.Errorf("не получилось распарсить conflict_Error short_url %w", err)
				}

				// Разрешаем ссылку относительно базового URL
				resolvedURL := baseURL.ResolveReference(shortURL)

				// Возвращаем окончательный URL и ошибку конфликта
				return resolvedURL.String(), ErrConflict
			}
		}

		if err == nil {
			// Если URL успешно сохранен, выходим из цикла
			break
		}
	}

	if err != nil {
		return "", fmt.Errorf("service can not save URL %w", err)
	}

	// Формируем окончательный короткий URL
	idPath, err := url.Parse(id)
	if err != nil {
		return "", fmt.Errorf("ошибка при разборе пути id: %w", err)
	}

	shortURL := base.ResolveReference(idPath)

	return shortURL.String(), nil
}

// NewShortURLService создает новый экземпляр ShortURLService с заданным количеством попыток генерации.
func NewShortURLService() *ShortURLService {
	return &ShortURLService{
		maxIteration: 10, //nolint:gomnd // Количество попыток генерации короткого URL
	}
}

// isConflictError проверяет, является ли ошибка ошибкой конфликта.
func isConflictError(err error) (storageErrors.ConflictError, bool) {
	var conflictErr storageErrors.ConflictError
	if errors.As(err, &conflictErr) {
		return conflictErr, true
	}
	return conflictErr, false
}
