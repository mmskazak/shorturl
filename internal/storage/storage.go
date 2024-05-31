package storage

import (
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage/infile"
	"mmskazak/shorturl/internal/storage/inmemory"
	"mmskazak/shorturl/internal/storage/postgresql"
)

type Storage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

func NewStorage(cfg *config.Config) (Storage, error) {
	switch {
	case cfg.DataBaseDSN != "":
		pg, err := postgresql.NewPostgreSQL(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
		}
		return pg, nil
	case cfg.FileStoragePath == "":
		sm, err := inmemory.NewInMemory()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize an in-memory store: %w", err)
		}
		return sm, nil
	case cfg.FileStoragePath != "":
		sf, err := infile.NewInFile(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize an in-file store: %w", err)
		}
		return sf, nil
	default:
		return nil, errors.New("error creating new storage")
	}
}
