package factory

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage"
	"mmskazak/shorturl/internal/storage/postgresql"

	"go.uber.org/zap"
)

func NewStorage(ctx context.Context, cfg *config.Config, zapLog *zap.SugaredLogger) (storage.Storage, error) {
	pqLogger := zapLog.With(zap.String("storage", "PostgreSQL"))
	pg, err := postgresql.NewPostgreSQL(ctx, cfg, pqLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
	}
	return pg, nil
}
