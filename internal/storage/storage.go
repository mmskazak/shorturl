package storage

import (
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage/infile"
	"mmskazak/shorturl/internal/storage/inmemory"
)

type Storage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

func NewStorage(storageType string, cfg *config.Config) (Storage, error) {
	switch storageType {
	case "inmemory":
		return inmemory.NewInMemory() //nolint:wrapcheck //ошибка обрабатывается далее
	case "infile":
		return infile.NewInFile(cfg) //nolint:wrapcheck //ошибка обрабатываетсяч далее
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}
