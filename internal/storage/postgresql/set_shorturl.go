package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// SetShortURL error:
// different error
// ErrKeyAlreadyExists
// ConflictError (ErrOriginalURLAlreadyExists).
func (s *PostgreSQL) SetShortURL(ctx context.Context, shortURL string, targetURL string, userID string) error {
	// Начало транзакции
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				s.zapLog.Infof("error rollback transaction: %v", errRollback)
			}
		}
	}()

	// Выполняем команду INSERT в контексте транзакции
	_, err = tx.Exec(ctx, `
        INSERT INTO urls (short_url, original_url, user_id)
        VALUES ($1, $2, $3)
    `, shortURL, targetURL, userID)

	if err != nil {
		return s.handleError(ctx, err, targetURL)
	}

	if err = tx.Commit(ctx); err != nil {
		s.zapLog.Infof("error committing transaction: %v", err)
	}

	// Если все успешно, err остается nil и транзакция будет зафиксирована
	return nil
}

func (s *PostgreSQL) handleError(ctx context.Context, err error, targetURL string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		switch pgErr.ConstraintName {
		case "unique_short_url":
			return storageErrors.ErrKeyAlreadyExists
		case "unique_original_url":
			var shortURL string
			err := s.pool.QueryRow(ctx, "SELECT short_url FROM urls WHERE original_url = $1", targetURL).Scan(&shortURL)
			if err != nil {
				return fmt.Errorf("error recive short URL by original: %w", err)
			}
			return storageErrors.ConflictError{
				ShortURL: shortURL,
				Err:      storageErrors.ErrOriginalURLAlreadyExists,
			}
		}
	}
	return fmt.Errorf("failed to insert record: %w", err)
}
