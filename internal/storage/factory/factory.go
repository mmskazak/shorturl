package factory

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage"
	"mmskazak/shorturl/internal/storage/infile"
	"mmskazak/shorturl/internal/storage/inmemory"
	"mmskazak/shorturl/internal/storage/postgresql"

	"go.uber.org/zap"
)

func NewStorage(ctx context.Context, cfg *config.Config, zapLog *zap.SugaredLogger) (storage.Storage, error) {
	switch {
	case cfg.DataBaseDSN != "":
		PostgreSQL := zapLog.With(zap.String("storage", "PostgreSQL"))
		pg, err := postgresql.NewPostgreSQL(ctx, cfg, PostgreSQL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
		}
		return pg, nil
	case cfg.FileStoragePath == "":
		inMemoryLog := zapLog.With(zap.String("storage", "InMemory"))
		sm, err := inmemory.NewInMemory(inMemoryLog)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize an in-memory store: %w", err)
		}
		return sm, nil
	default:
		InFile := zapLog.With(zap.String("storage", "InFile"))
		sf, err := infile.NewInFile(ctx, cfg, InFile)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize an in-file store: %w", err)
		}
		return sf, nil
	}
}
