package shorturlservice

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

var ErrServiceGenerateID = errors.New("generateID failed")
var ErrConflict = errors.New("error original url already exist")

type IGenIDForURL interface {
	Generate() (string, error)
}

type DTOShortURL struct {
	OriginalURL string
	UserID      string
	BaseHost    string
	Deleted     bool
}

type ShortURLService struct {
	maxIteration int
}

func (s *ShortURLService) GenerateShortURL(
	ctx context.Context,
	dto DTOShortURL,
	generator IGenIDForURL,
	data storage.Storage,
) (string, error) {
	var err error
	base, err := url.Parse(dto.BaseHost)
	if err != nil {
		return "", fmt.Errorf("ошибка при разборе базового URL: %w", err)
	}

	var id string
	for range s.maxIteration {
		id, err = generator.Generate()
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrServiceGenerateID, err)
		}

		err = data.SetShortURL(ctx, id, dto.OriginalURL, dto.UserID, dto.Deleted)
		if err != nil {
			conflictError, ok := isConflictError(err)
			if ok {
				// Парсинг базового хоста
				baseURL, err := url.Parse(dto.BaseHost)
				if err != nil {
					return "", fmt.Errorf("не получилось распарсить dto base_host %w", err)
				}

				// Парсинг сокращенного URL
				shortURL, err := url.Parse(conflictError.ShortURL)
				if err != nil {
					return "", fmt.Errorf("не получилось распарсить conflict_Error short_url %w", err)
				}

				// Разрешение ссылки относительно базового URL-адреса
				resolvedURL := baseURL.ResolveReference(shortURL)

				// Получение строки URL
				finalURL := resolvedURL.String()

				return finalURL, ErrConflict
			}
		}

		if err == nil {
			break
		}
	}

	if err != nil {
		return "", fmt.Errorf("service can not save URL %w", err)
	}

	idPath, err := url.Parse(id)
	if err != nil {
		return "", fmt.Errorf("ошибка при разборе пути id: %w", err)
	}

	shortURL := base.ResolveReference(idPath)

	return shortURL.String(), nil
}

func NewShortURLService() *ShortURLService {
	return &ShortURLService{
		maxIteration: 10, //nolint:gomnd //количество попыток генерирования short url
	}
}

func isConflictError(err error) (storageErrors.ConflictError, bool) {
	var conflictErr storageErrors.ConflictError
	if errors.As(err, &conflictErr) {
		return conflictErr, true
	}
	return conflictErr, false
}
