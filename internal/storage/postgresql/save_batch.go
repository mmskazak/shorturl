package postgresql

import (
	"context"
	"database/sql"
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
	const maxRetries = 10 // Максимальное количество попыток для каждой записи

	lenItems := len(items)
	if lenItems == 0 {
		return nil, errors.New("batch with original URL is empty")
	}

	outputs := make([]storage.Output, 0, lenItems)

	// Начало транзакции
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				s.zapLog.Warnf("error rolling back transaction: %v", err)
			}
		}
	}()

	incomingMap := make(map[string]string)
	for i := range lenItems {
		item := items[i]
		incomingMap[item.OriginalURL] = item.CorrelationID

		// Повторные попытки вставки с разной короткой URL
		for retry := range maxRetries {
			idShortURL, err := generator.Generate()
			if err != nil {
				return nil, fmt.Errorf("error generating ID for URL: %w", err)
			}

			stmt := "INSERT INTO urls(short_url, original_url, user_id) VALUES ($1, $2, $3) RETURNING short_url, original_url"
			batch := &pgx.Batch{}
			batch.Queue(stmt, idShortURL, item.OriginalURL, userID)

			// Отправляем батчевый запрос и получаем результаты
			batchResults := tx.SendBatch(ctx, batch)
			var shortURL string
			var originalURL string
			err = batchResults.QueryRow().Scan(&shortURL, &originalURL)

			var pgErr *pgconn.PgError
			if err != nil {
				if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
					log.Printf("Unique constraint violation on retry %d for URL %s: %v", retry, item.OriginalURL, err)
					// Если ошибка уникальности, продолжаем попытки
					if pgErr.ConstraintName == "unique_short_url" {
						// Конфликт по уникальному короткому URL, генерируем заново и пробуем ещё раз
						continue
					}
					if pgErr.ConstraintName == "unique_original_url" {
						return nil, storageErrors.ErrUniqueViolation
					}
				} else {
					return nil, fmt.Errorf("error inserting data: %w", err)
				}
			}

			// Если вставка успешна, завершаем цикл повторных попыток
			fullShortURL, err := storage.GetFullShortURL(baseHost, shortURL)
			if err != nil {
				return nil, fmt.Errorf("error getting full short URL: %w", err)
			}

			outputs = append(outputs, storage.Output{
				CorrelationID: incomingMap[originalURL],
				ShortURL:      fullShortURL,
			})

			if err := batchResults.Close(); err != nil {
				return nil, fmt.Errorf("error closing batch results: %w", err)
			}

			// Успешная вставка, выходим из цикла повторных попыток
			break
		}

		// Если все попытки завершились неудачей, возвращаем ошибку
		if i == lenItems-1 {
			return nil, fmt.Errorf("failed to insert URL %s after %d retries", item.OriginalURL, maxRetries)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("error committing transaction: %v", err)
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return outputs, nil
}
