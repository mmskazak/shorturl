package factory

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/contracts"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage/infile"
	"mmskazak/shorturl/internal/storage/inmemory"
	"mmskazak/shorturl/internal/storage/postgresql"

	"go.uber.org/zap"
)

// NewStorage создает и возвращает экземпляр хранилища данных в зависимости от конфигурации.
// Функция принимает контекст, конфигурацию и логгер, и возвращает интерфейс хранилища и ошибку, если она возникла.
//
// В зависимости от настроек конфигурации выбирается один из трех типов хранилищ:
// 1. PostgreSQL - если указан DSN для базы данных PostgreSQL.
// 2. InMemory - если путь к файловому хранилищу не указан, создается хранилище в памяти.
// 3. InFile - если указано, но DSN не задан, создается файловое хранилище.
//
// Если инициализация хранилища завершается ошибкой, функция возвращает ошибку с пояснением.
func NewStorage(ctx context.Context, cfg *config.Config, zapLog *zap.SugaredLogger) (contracts.Storage, error) {
	switch {
	// Если задан DSN для базы данных PostgreSQL, создаем и возвращаем соединение с PostgreSQL.
	case cfg.DataBaseDSN != "":
		pqLogger := zapLog.With(zap.String("storage", "PostgreSQL"))
		pg, err := postgresql.NewPostgreSQL(ctx, cfg, pqLogger)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
		}
		return pg, nil

	// Если путь к файловому хранилищу не задан, создаем и возвращаем хранилище в памяти.
	case cfg.FileStoragePath == "":
		inMemoryLogger := zapLog.With(zap.String("storage", "InMemory"))
		sm, err := inmemory.NewInMemory(inMemoryLogger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize an in-memory store: %w", err)
		}
		return sm, nil

	// В остальных случаях создаем и возвращаем файловое хранилище.
	default:
		inFileLogger := zapLog.With(zap.String("storage", "InFile"))
		sf, err := infile.NewInFile(ctx, cfg, inFileLogger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize an in-file store: %w", err)
		}
		return sf, nil
	}
}
