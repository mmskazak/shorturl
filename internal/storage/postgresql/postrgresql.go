package postgresql

import (
	"context"
	"fmt"

	"mmskazak/shorturl/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgreSQL struct {
	pool   *pgxpool.Pool
	zapLog *zap.SugaredLogger
}

func NewPostgreSQL(ctx context.Context, cfg *config.Config, zapLog *zap.SugaredLogger) (*PostgreSQL, error) {
	pool, err := pgxpool.New(ctx, cfg.DataBaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to dbshorturl: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping dbshorturl connection: %w", err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_url VARCHAR(255) NOT NULL,
			original_url TEXT NOT NULL,
		    user_id VARCHAR(255),
			CONSTRAINT unique_short_url UNIQUE (short_url),
			CONSTRAINT unique_original_url UNIQUE (original_url)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table shorturl: %w", err)
	}

	return &PostgreSQL{
		pool:   pool,
		zapLog: zapLog,
	}, nil
}

func (p *PostgreSQL) Ping(ctx context.Context) error {
	err := p.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	return nil
}

func (p *PostgreSQL) Close() error {
	if p.pool == nil {
		return nil
	}
	p.pool.Close()
	return nil
}
