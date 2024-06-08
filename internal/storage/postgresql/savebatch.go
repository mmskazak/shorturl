package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// SaveBatch inserts a batch of short URL mappings into the database.
func (p *PostgreSQL) SaveBatch(
	ctx context.Context,
	items []storage.Incoming,
	baseHost string,
) ([]storage.Output, error) {
	lenItems := len(items)
	if lenItems == 0 {
		return nil, errors.New("batch with original URL is empty")
	}

	outputs := make([]storage.Output, 0, lenItems)

	// Начало транзакции
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				log.Printf("error rolling back transaction: %v", errRollback)
			}
		} else {
			if err = tx.Commit(ctx); err != nil {
				log.Printf("error committing transaction: %v", err)
			}
		}
	}()

	// Используем pgx.Batch для отправки множества команд в одном запросе
	batch := &pgx.Batch{}

	for _, item := range items {
		stmt := "INSERT INTO urls(short_url, original_url) VALUES ($1, $2) RETURNING short_url"
		batch.Queue(stmt, item.CorrelationID, item.OriginalURL)
	}

	// Отправляем батчевый запрос и получаем результаты
	batchResults := tx.SendBatch(ctx, batch)
	defer func(batchResults pgx.BatchResults) {
		err := batchResults.Close()
		if err != nil {
			log.Printf("error closing batch results: %v", err)
		}
	}(batchResults)

	for range lenItems {
		var shortURL string
		err = batchResults.QueryRow().Scan(&shortURL)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == ErrCodeDatabaseUniqueViolation {
				return nil, storageErrors.ErrUniqueViolation
			}
			return nil, fmt.Errorf("error inserting data: %w", err)
		}

		fullShortURL, err := storage.GetFullShortURL(baseHost, shortURL)
		if err != nil {
			return nil, fmt.Errorf("error getting full short URL: %w", err)
		}

		outputs = append(outputs, storage.Output{
			CorrelationID: shortURL,
			ShortURL:      fullShortURL,
		})
	}

	if err := batchResults.Close(); err != nil {
		return nil, fmt.Errorf("error closing batch results: %w", err)
	}

	return outputs, nil
}
