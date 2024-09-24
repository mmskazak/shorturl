package postgresql

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage/postgresql/interfaces"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// PostgreSQL представляет собой структуру для работы с PostgreSQL.
type PostgreSQL struct {
	pool   interfaces.Database
	zapLog *zap.SugaredLogger
}

//go:embed migrations/*
var embedMigrations embed.FS

const migrationsDir = "migrations"

// NewPostgreSQL инициализирует новое соединение с PostgreSQL, проверяет подключение,
// запускает миграции и запускает горутину для периодического удаления помеченных записей.
func NewPostgreSQL(ctx context.Context, cfg *config.Config, zapLog *zap.SugaredLogger) (*PostgreSQL, error) {
	zapLog.Infof("initializing PostgreSQL")
	pool, err := pgxpool.New(ctx, cfg.DataBaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to dbshorturl: %w", err)
	}

	zapLog.Infof("ping dbshorturl")
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping dbshorturl connection: %w", err)
	}
	zapLog.Infof("run migrations")
	if err := runMigrations(cfg.DataBaseDSN, zapLog); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	go hardDeleteSoftDeletedURLs(ctx, pool, zapLog)

	return &PostgreSQL{
		pool:   pool,
		zapLog: zapLog,
	}, nil
}

// Ping проверяет доступность соединения с PostgreSQL.
func (s *PostgreSQL) Ping(ctx context.Context) error {
	err := s.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	return nil
}

// Close закрывает пул соединений с PostgreSQL.
func (s *PostgreSQL) Close() error {
	if s.pool == nil {
		return errors.New("postgreSQL pool is nil")
	}
	s.pool.Close()
	return nil
}

// runMigrations запускает миграции базы данных из встраиваемой файловой системы.
func runMigrations(dsn string, zapLog *zap.SugaredLogger) error {
	zapLog.Infof("Путь к директории миграций: %s", migrationsDir)
	dir, err := iofs.New(embedMigrations, migrationsDir)
	if err != nil {
		return fmt.Errorf("error opening migrations directory: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", dir, dsn)
	if err != nil {
		return fmt.Errorf("error opening migrations directory: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("error running migrations: %w", err)
		}
	}
	return nil
}
