package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"log"

	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const batchSize = 5000

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
		if err = tx.Rollback(ctx); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("error rollback transaction: %v", err)
			}
		}
	}()

	// Используем pgx.Batch для отправки множества команд в одном запросе
	batch := &pgx.Batch{}

	for i := range lenItems {
		item := items[i]
		stmt := "INSERT INTO urls(short_url, original_url) VALUES ($1, $2) RETURNING short_url"
		batch.Queue(stmt, item.CorrelationID, item.OriginalURL)

		// Если количество запросов в батче достигло предела или это последний элемент,
		// то отправляем батчевый запрос и обрабатываем результаты
		if (i+1)%batchSize == 0 || i == lenItems-1 {
			// Отправляем батчевый запрос и получаем результаты
			batchResults := tx.SendBatch(ctx, batch)

			for range batch.Len() {
				var shortURL string
				err = batchResults.QueryRow().Scan(&shortURL)
				if err != nil {
					var pgErr *pgconn.PgError
					if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
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

			// Очищаем батч для следующей порции запросов
			batch = &pgx.Batch{}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("error committing transaction: %v", err)
	}

	return outputs, nil
}
