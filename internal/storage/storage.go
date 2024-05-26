package storage

import (
	"errors"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage/infile"
	"mmskazak/shorturl/internal/storage/inmemory"
)

type Storage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

func NewStorage(cfg *config.Config) (Storage, error) {
	switch {
	case cfg.FileStoragePath == "":
		return inmemory.NewInMemory() //nolint:wrapcheck //ошибка обрабатывается далее
	case cfg.FileStoragePath != "":
		return infile.NewInFile(cfg) //nolint:wrapcheck //ошибка обрабатываетсяч далее
	default:
		return nil, errors.New("error creating new storage")
	}
}
