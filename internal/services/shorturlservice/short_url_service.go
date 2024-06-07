package shorturlservice

import (
	"errors"
	"fmt"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"net/url"
)

var ErrOriginalURLIsEmpty = errors.New("originalURL is empty")
var ErrBaseHostIsEmpty = errors.New("base host is empty")
var ErrServiceGenerateID = errors.New("generateID failed")
var ErrConflict = errors.New("status 409 conflict")

type IGenIDForURL interface {
	Generate(int) (string, error)
}

type Storage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

type DTOShortURL struct {
	OriginalURL  string
	BaseHost     string
	MaxIteration int
	LengthID     int
}

type ShortURLService struct{}

func (s *ShortURLService) GenerateShortURL(dto DTOShortURL, generator IGenIDForURL, storage Storage) (string, error) {
	if dto.OriginalURL == "" {
		return "", ErrOriginalURLIsEmpty
	}

	if dto.BaseHost == "" {
		return "", ErrBaseHostIsEmpty
	}

	var err error
	base, err := url.Parse(dto.BaseHost)
	if err != nil {
		return "", fmt.Errorf("ошибка при разборе базового URL: %w", err)
	}

	var id string
	for range dto.MaxIteration {
		id, err = generator.Generate(dto.LengthID)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrServiceGenerateID, err)
		}

		err = storage.SetShortURL(id, dto.OriginalURL)
		if err != nil {
			conflictError, ok := IsConflictError(err)
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
		return "", fmt.Errorf("ошибка при разборе пути ID: %w", err)
	}

	shortURL := base.ResolveReference(idPath)

	return shortURL.String(), nil
}

func NewShortURLService() *ShortURLService {
	return &ShortURLService{}
}

func IsConflictError(err error) (*storageErrors.ConflictError, bool) {
	var conflictErr *storageErrors.ConflictError
	if errors.As(err, &conflictErr) {
		return conflictErr, true
	}
	return nil, false
}
