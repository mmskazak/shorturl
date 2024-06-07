package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"github.com/jackc/pgx/v5/pgconn"

	"mmskazak/shorturl/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const ErrDatabaseUniqueViolation = "23505"

type PostgreSQL struct {
	db *sql.DB
}

func NewPostgreSQL(cfg *config.Config) (*PostgreSQL, error) {
	// Подключаемся к базе данных dbshorturl
	dbShortURL, err := sql.Open("pgx", cfg.DataBaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to dbshorturl: %w", err)
	}

	// Проверяем соединение с dbshorturl
	err = dbShortURL.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping dbshorturl connection: %w", err)
	}

	// Создаем таблицу shorturl, если она не существует
	_, err = dbShortURL.Exec(`
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
		db: dbShortURL,
	}, nil
}

func (p *PostgreSQL) GetShortURL(shortURL string) (string, error) {
	var originalURL string
	// Выполняем запрос SQL для получения original_url по short_url
	err := p.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("short URL not found %w", err)
		}
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}
	return originalURL, nil
}

func (p *PostgreSQL) SetShortURL(shortURL string, targetURL string) error {
	_, err := p.db.Exec(`
        INSERT INTO urls (short_url, original_url)
        VALUES ($1, $2)
    `, shortURL, targetURL)

	if err != nil {
		err := p.duplicateOriginalURL(err, targetURL)
		if err != nil {
			return fmt.Errorf("duplicate origin URL: %w", err)
		}

		err = p.duplicateShortURL(err, targetURL)
		if err != nil {
			return fmt.Errorf("duplicate short URL: %w", err)
		}
		return fmt.Errorf("failed to insert record: %w", err)
	}
	return nil
}

func (p *PostgreSQL) Ping() error {
	err := p.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	return nil
}

func (p *PostgreSQL) Close() error {
	if p.db == nil {
		return nil
	}

	err := p.db.Close()
	if err != nil {
		return fmt.Errorf("error closing database connection: %w", err)
	}

	return nil
}

func (p *PostgreSQL) duplicateOriginalURL(err error, targetURL string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == ErrDatabaseUniqueViolation {
		var shortURL string
		errQuery := p.db.QueryRow(`
                SELECT short_url FROM urls WHERE original_url = $1
            `, targetURL).Scan(&shortURL)
		if errQuery != nil {
			return fmt.Errorf("falied to get short url for original url: %w", err)
		}

		if shortURL != "" {
			return &storageErrors.ConflictError{
				ShortURL: shortURL,
				Err:      err,
			}
		}
	}
	return nil
}

func (p *PostgreSQL) duplicateShortURL(err error, shortURL string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == ErrDatabaseUniqueViolation {
		var originalURL string
		errQuery := p.db.QueryRow(`
                SELECT original_url FROM urls WHERE short_url = $1
            `, shortURL).Scan(&originalURL)
		if errQuery != nil {
			return fmt.Errorf("falied to get short url for original url: %w", err)
		}

		if originalURL != "" {
			return storageErrors.ErrKeyAlreadyExists
		}
	}
	return nil
}
