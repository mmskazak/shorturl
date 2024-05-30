package shorturlservice

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrOriginalURLIsEmpty = errors.New("originalURL is empty")
var ErrBaseHostIsEmpty = errors.New("base host is empty")
var ErrServiceGenerateID = errors.New("generateID failed")

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
	var id string
	for range dto.MaxIteration {
		id, err = generator.Generate(dto.LengthID)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrServiceGenerateID, err)
		}

		err = storage.SetShortURL(id, dto.OriginalURL)
		if err == nil {
			break
		}
	}

	if err != nil {
		return "", errors.New("service can not save URL")
	}

	base, err := url.Parse(dto.BaseHost)
	if err != nil {
		return "", fmt.Errorf("ошибка при разборе базового URL: %w", err)
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
