package postgresql

import (
	"context"
	"errors"
	"fmt"

	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/models"

	"github.com/jackc/pgerrcode"

	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const batchSize = 5000

// SaveBatch сохраняет пакет URL-адресов в базе данных PostgreSQL.
// Возвращает список сохраненных коротких URL-адресов или ошибку в случае неудачи.
//
// Ошибки:
// - ErrKeyAlreadyExists: короткий URL уже существует
// - ConflictError: оригинальный URL уже существует.
func (s *PostgreSQL) SaveBatch(
	ctx context.Context,
	items []models.Incoming,
	baseHost string,
	userID string,
	generator contracts.IGenIDForURL,
) ([]models.Output, error) {
	lenItems := len(items)
	if lenItems == 0 {
		return nil, errors.New("batch with original URL is empty")
	}

	outputs := make([]models.Output, 0, lenItems)

	// Начало транзакции
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			s.zapLog.Errorf("error rollback transaction: %v", err)
		}
	}()

	// Используем pgx.Batch для отправки множества команд в одном запросе
	batch := &pgx.Batch{}

	incomingMap := make(map[string]string)
	for i := range lenItems {
		item := items[i]
		incomingMap[item.OriginalURL] = item.CorrelationID

		idShortURL, err := generator.Generate()
		if err != nil {
			return nil, fmt.Errorf("error generating ID for URL: %w", err)
		}
		stmt := "INSERT INTO urls(short_url, original_url, user_id, deleted) VALUES ($1, $2, $3, $4) " +
			"RETURNING short_url, original_url"
		batch.Queue(stmt, idShortURL, item.OriginalURL, userID, false)

		// Если количество запросов в батче достигло предела или это последний элемент,
		// то отправляем батчевый запрос и обрабатываем результаты
		if (i+1)%batchSize == 0 || i == lenItems-1 {
			// Отправляем батчевый запрос и получаем результаты
			batchResults := tx.SendBatch(ctx, batch)

			for range batch.Len() {
				var shortURL string
				var originalURL string
				err = batchResults.QueryRow().Scan(&shortURL, &originalURL)
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
					switch pgErr.ConstraintName {
					case "unique_short_url":
						return nil, storageErrors.ErrKeyAlreadyExists
					case "unique_original_url":
						return nil, storageErrors.ErrUniqueViolation
					}
				}
				if err != nil {
					return nil, fmt.Errorf("error inserting data: %w", err)
				}

				fullShortURL, err := storage.GetFullShortURL(baseHost, shortURL)
				if err != nil {
					return nil, fmt.Errorf("error getting full short URL: %w", err)
				}

				outputs = append(outputs, models.Output{
					CorrelationID: incomingMap[originalURL],
					ShortURL:      fullShortURL,
				})
			}

			if err := batchResults.Close(); err != nil {
				return nil, fmt.Errorf("error closing batch results: %w", err)
			}

			// Очищаем батч для следующей порции запросов
			batch = &pgx.Batch{}
			// Очищаем вспомогательную структуру
			incomingMap = make(map[string]string)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		s.zapLog.Warnf("error committing batch transaction: %v", err)
	}

	return outputs, nil
}
