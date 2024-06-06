package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func (p *PostgreSQL) SaveBatch(items []storage.Incoming, baseHost string) ([]storage.Output, error) {
	lenItems := len(items)

	if lenItems == 0 {
		return nil, errors.New("batch with original URL is empty")
	}

	outputs := make([]storage.Output, 0, lenItems)
	tx, err := p.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Printf("error rolling back transaction: %v", errRollback)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("error committing transaction: %v", err)
			}
		}
	}()

	batchSize := 4
	// Переменная до цикла, чтобы потом вызвать defer, а не в цикле
	var rows *sql.Rows
	for start := 0; start < lenItems; start += batchSize {
		end := start + batchSize
		if end > lenItems {
			end = lenItems
		}

		// Генерация SQL-запроса для текущей партии
		stmt := generateURLsStatement(end - start)

		// Подготовка значений для вставки
		args := make([]interface{}, 0, (end-start)*2) //nolint:gomnd // потому что 2 элемента в каждом insert
		for _, item := range items[start:end] {
			args = append(args, item.CorrelationID, item.OriginalURL)
		}

		rows, err = tx.Query(stmt, args...)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == ErrDatabaseUniqueViolation {
				return nil, storageErrors.ErrOriginalURLAlreadyExists
			}
			return nil, fmt.Errorf("error inserting data: %w", err)
		}

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
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}(rows)

	return outputs, nil
}

func generateURLsStatement(count int) string {
	const stmtTmpl = `INSERT INTO urls(short_url, original_url) VALUES %s RETURNING short_url`

	valuesParts := make([]string, 0, count)
	for i := range count {
		valuesParts = append(valuesParts,
			fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)) //nolint:gomnd //собираем переменые под вставку данных
	}

	return fmt.Sprintf(stmtTmpl, strings.Join(valuesParts, ","))
}
