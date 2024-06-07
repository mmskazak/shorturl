package postgresql

import (
	"context"
	"errors"
	"fmt"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"mmskazak/shorturl/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const ErrDatabaseUniqueViolation = "23505"

type PostgreSQL struct {
	conn *pgx.Conn
}

// NewPostgreSQL initializes a new PostgreSQL connection using pgx.
func NewPostgreSQL(ctx context.Context, cfg *config.Config) (*PostgreSQL, error) {
	// Подключаемся к базе данных dbshorturl
	conn, err := pgx.Connect(ctx, cfg.DataBaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to dbshorturl: %w", err)
	}

	// Проверяем соединение с dbshorturl
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping dbshorturl connection: %w", err)
	}

	// Создаем таблицу shorturl, если она не существует
	_, err = conn.Exec(ctx, `
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
		conn: conn,
	}, nil
}

// GetShortURL retrieves the original URL for the given short URL.
func (p *PostgreSQL) GetShortURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	err := p.conn.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("short URL not found: %w", err)
		}
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}
	return originalURL, nil
}

// SetShortURL inserts a new short URL and its corresponding original URL.
func (p *PostgreSQL) SetShortURL(ctx context.Context, shortURL string, targetURL string) error {
	_, err := p.conn.Exec(ctx, `
		INSERT INTO urls (short_url, original_url)
		VALUES ($1, $2)
	`, shortURL, targetURL)

	if err != nil {
		if err := p.handleDuplicateError(ctx, err, shortURL, targetURL); err != nil {
			return err
		}
		return fmt.Errorf("failed to insert record: %w", err)
	}
	return nil
}

// Ping checks the connection to the database.
func (p *PostgreSQL) Ping(ctx context.Context) error {
	err := p.conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	return nil
}

// Close closes the connection to the database.
func (p *PostgreSQL) Close(ctx context.Context) error {
	if p.conn == nil {
		return nil
	}
	err := p.conn.Close(ctx)
	if err != nil {
		return fmt.Errorf("error closing database connection: %w", err)
	}
	return nil
}

// handleDuplicateError checks and handles unique constraint violations for the short and original URLs.
func (p *PostgreSQL) handleDuplicateError(ctx context.Context, err error, shortURL string, targetURL string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == ErrDatabaseUniqueViolation {
		if dupeErr := p.checkDuplicateOriginalURL(ctx, targetURL); dupeErr != nil {
			return dupeErr
		}
		if dupeErr := p.checkDuplicateShortURL(ctx, shortURL); dupeErr != nil {
			return dupeErr
		}
	}
	return nil
}

// checkDuplicateOriginalURL checks if the original URL already exists and returns the corresponding short URL.
func (p *PostgreSQL) checkDuplicateOriginalURL(ctx context.Context, targetURL string) error {
	var shortURL string
	err := p.conn.QueryRow(ctx, `
		SELECT short_url FROM urls WHERE original_url = $1
	`, targetURL).Scan(&shortURL)
	if err != nil {
		return fmt.Errorf("failed to get short URL for original URL: %w", err)
	}
	if shortURL != "" {
		return &storageErrors.ConflictError{
			ShortURL: shortURL,
			Err:      fmt.Errorf("original URL already exists %w", err),
		}
	}
	return nil
}

// checkDuplicateShortURL checks if the short URL already exists.
func (p *PostgreSQL) checkDuplicateShortURL(ctx context.Context, shortURL string) error {
	var originalURL string
	err := p.conn.QueryRow(ctx, `
		SELECT original_url FROM urls WHERE short_url = $1
	`, shortURL).Scan(&originalURL)
	if err != nil {
		return fmt.Errorf("failed to get original URL for short URL: %w", err)
	}
	if originalURL != "" {
		return storageErrors.ErrKeyAlreadyExists
	}
	return nil
}
