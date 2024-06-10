package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgerrcode"

	"mmskazak/shorturl/internal/config"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQL struct {
	pool *pgxpool.Pool
}

func NewPostgreSQL(ctx context.Context, cfg *config.Config) (*PostgreSQL, error) {
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
			CONSTRAINT unique_short_url UNIQUE (short_url),
			CONSTRAINT unique_original_url UNIQUE (original_url)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table shorturl: %w", err)
	}

	return &PostgreSQL{
		pool: pool,
	}, nil
}

func (p *PostgreSQL) GetShortURL(ctx context.Context, shortURL string) (string, error) {

	var originalURL string
	err := p.pool.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("short URL not found: %w", storageErrors.ErrNotFound)
		}
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}
	return originalURL, nil
}

func (p *PostgreSQL) SetShortURL(ctx context.Context, shortURL string, targetURL string) error {
	// Начало транзакции
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("error rollback transaction: %v", errRollback)
			}
		}
	}()

	// Выполняем команду INSERT в контексте транзакции
	_, err = tx.Exec(ctx, `
        INSERT INTO urls (short_url, original_url)
        VALUES ($1, $2)
    `, shortURL, targetURL)

	if err != nil {
		return p.handleError(ctx, err, targetURL)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("error committing transaction: %v", err)
	}

	// Если все успешно, err остается nil и транзакция будет зафиксирована
	return nil
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

func (p *PostgreSQL) handleError(ctx context.Context, err error, targetURL string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		switch pgErr.ConstraintName {
		case "unique_short_url":
			return storageErrors.ErrKeyAlreadyExists
		case "unique_original_url":
			var shortURL string
			err := p.pool.QueryRow(ctx, "SELECT short_url FROM urls WHERE original_url = $1", targetURL).Scan(&shortURL)
			if err != nil {
				return fmt.Errorf("error recive short URL by original: %w", err)
			}
			return &storageErrors.ConflictError{
				ShortURL: shortURL,
				Err:      storageErrors.ErrOriginalURLAlreadyExists,
			}
		}
	}
	return fmt.Errorf("failed to insert record: %w", err)
}
