package postgresql

import (
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

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

	stmt, err := tx.Prepare("INSERT INTO urls (short_url, original_url) VALUES ($1, $2) RETURNING id, short_url")
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return nil, fmt.Errorf("error rolback transaction %w", err)
		}
		return nil, fmt.Errorf("ошибка подготовки оператора: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("ошибка при закрытии stmt: %v", err)
		}
	}()

	for _, item := range items {
		_, err = stmt.Exec(&item.CorrelationID, &item.OriginalURL)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return nil, fmt.Errorf("error rolback transaction %w", errRollback)
			}

			// Проверим, является ли ошибка нарушением ограничения внешнего ключа
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == ErrDatabaseUniqueViolation {
				return nil, storageErrors.ErrOriginalURLAlreadyExists
			}

			return nil, fmt.Errorf("error inserting data: %w", err)
		}

		fullShortURL, err := storage.GetFullShortURL(baseHost, item.CorrelationID)
		if err != nil {
			return nil, fmt.Errorf("error getFullShortURL from two parts %w", err)
		}

		outputs = append(outputs, storage.Output{
			CorrelationID: item.CorrelationID,
			ShortURL:      fullShortURL,
		})
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return outputs, nil
}
