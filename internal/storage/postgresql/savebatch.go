package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// SaveBatch inserts a batch of short URL mappings into the database.
func (p *PostgreSQL) SaveBatch(
	ctx context.Context,
	items []storage.Incoming,
	baseHost string) ([]storage.Output, error) {
	lenItems := len(items)
	if lenItems == 0 {
		return nil, errors.New("batch with original URL is empty")
	}

	outputs := make([]storage.Output, 0, lenItems)

	// Начало транзакции
	tx, err := p.conn.Begin(ctx)
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

	batchSize := 4
	for start := 0; start < lenItems; start += batchSize {
		end := start + batchSize
		if end > lenItems {
			end = lenItems
		}

		// Генерация SQL-запроса для текущей партии
		stmt := generateURLsStatement(end - start)

		// Подготовка значений для вставки
		args := make([]interface{}, 0, (end-start)*2) //nolint:gomnd //на каждую запись дав значения
		for _, item := range items[start:end] {
			args = append(args, item.CorrelationID, item.OriginalURL)
		}

		// Выполнение запроса с возвращением коротких URL
		rows, err := tx.Query(ctx, stmt, args...)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == ErrDatabaseUniqueViolation {
				return nil, storageErrors.ErrOriginalURLAlreadyExists
			}
			return nil, fmt.Errorf("error inserting data: %w", err)
		}

		// Обработка результатов запроса
		for rows.Next() {
			var shortURL string
			if err := rows.Scan(&shortURL); err != nil {
				return nil, fmt.Errorf("error scanning row: %w", err)
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

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over rows: %w", err)
		}

		rows.Close()
	}

	return outputs, nil
}

// generateURLsStatement generates an SQL INSERT statement for the batch of URLs.
func generateURLsStatement(count int) string {
	const stmtTmpl = `INSERT INTO urls(short_url, original_url) VALUES %s RETURNING short_url`

	valuesParts := make([]string, count)
	for i := range count {
		valuesParts[i] = fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2) //nolint:gomnd //$1 $2 $3 $4
	}

	return fmt.Sprintf(stmtTmpl, strings.Join(valuesParts, ","))
}
