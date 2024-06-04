package postgresql

import (
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/storage"
	"net/url"
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
		if cerr := stmt.Close(); cerr != nil {
			log.Printf("ошибка при закрытии stmt: %v", cerr)
		}
	}()

	for _, item := range items {
		_, err = stmt.Exec(&item.CorrelationID, &item.OriginalURL)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return nil, fmt.Errorf("error rolback transaction %w", err)
			}
			return nil, fmt.Errorf("error inserting data: %w", err)
		}

		fullShortURL, err := getFullShortURL(baseHost, item.CorrelationID)
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

func getFullShortURL(baseHost, correlationID string) (string, error) {
	u, err := url.Parse(baseHost)
	if err != nil {
		return "", fmt.Errorf("error parsing baseHost: %w", err)
	}
	// ResolveReference correctly concatenates the base URL and the path
	u = u.ResolveReference(&url.URL{Path: correlationID})
	return u.String(), nil
}
