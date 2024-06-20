package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgerrcode"

	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// SaveBatch error:.
// Different error
// ErrKeyAlreadyExists
// ConflictError (ErrOriginalURLAlreadyExists).
func (s *PostgreSQL) SaveBatch(
	ctx context.Context,
	items []storage.Incoming,
	baseHost string,
	userID string,
	generator storage.IGenIDForURL,
) ([]storage.Output, error) {
	if len(items) == 0 {
		return nil, errors.New("batch with original URLs is empty")
	}

	outputs := make([]storage.Output, 0, len(items))
	incomingMap := make(map[string]string)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.zapLog.Warnf("error rolling back transaction: %v", err)
		}
	}()

	stmt := "INSERT INTO urls(short_url, original_url, user_id, deleted) " +
		"VALUES ($1, $2, $3, $4) RETURNING short_url, original_url"
	batch := &pgx.Batch{}

	for _, item := range items {
		incomingMap[item.OriginalURL] = item.CorrelationID

		// Generate short URL
		idShortURL, err := generator.Generate()
		if err != nil {
			return nil, fmt.Errorf("error generating ID for URL: %w", err)
		}

		// Add to batch
		batch.Queue(stmt, idShortURL, item.OriginalURL, userID, false)
	}

	// Execute the batch
	batchResults := tx.SendBatch(ctx, batch)
	defer func() {
		if err := batchResults.Close(); err != nil {
			log.Printf("error closing batch results: %v", err)
		}
	}()

	for range items {
		var shortURL string
		var originalURL string

		// Get the result of the batch execution
		err := batchResults.QueryRow().Scan(&shortURL, &originalURL)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					switch pgErr.ConstraintName {
					case "unique_short_url":
						return nil, storageErrors.ErrKeyAlreadyExists
					case "unique_original_url":
						return nil, storageErrors.ErrUniqueViolation
					}
				}
			}
			return nil, fmt.Errorf("error inserting data: %w", err)
		}

		fullShortURL, err := storage.GetFullShortURL(baseHost, shortURL)
		if err != nil {
			return nil, fmt.Errorf("error getting full short URL: %w", err)
		}

		outputs = append(outputs, storage.Output{
			CorrelationID: incomingMap[originalURL],
			ShortURL:      fullShortURL,
		})
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("error committing transaction: %v", err)
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return outputs, nil
}
