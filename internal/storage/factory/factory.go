package factory

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage"
	"mmskazak/shorturl/internal/storage/infile"
	"mmskazak/shorturl/internal/storage/inmemory"
	"mmskazak/shorturl/internal/storage/postgresql"
)

func NewStorage(ctx context.Context, cfg *config.Config) (storage.Storage, error) {
	switch {
	case cfg.DataBaseDSN != "":
		pg, err := postgresql.NewPostgreSQL(ctx, cfg)
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
	default:
		sf, err := infile.NewInFile(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize an in-file store: %w", err)
		}
		return sf, nil
	}
}
