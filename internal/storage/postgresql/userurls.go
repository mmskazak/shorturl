package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/storage"
	"net/url"

	"github.com/jackc/pgx/v5"
)

func (p *PostgreSQL) GetUserURLs(ctx context.Context, userID string, baseHost string) ([]storage.URL, error) {
	// Определяем SQL-запрос для получения URL-адресов пользователя
	query := `
		SELECT short_url, original_url
		FROM urls
		WHERE user_id = $1
	`

	// Начало транзакции
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("error rolling back transaction: %v", err)
		}
	}()

	// Выполняем запрос и получаем строки
	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	// Создаем слайс для хранения результатов
	var urls []storage.URL

	// Обрабатываем результаты запроса
	for rows.Next() {
		var storageURL storage.URL
		if err := rows.Scan(&storageURL.ShortURL, &storageURL.OriginalURL); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Парсим базовый хост
		baseURL, err := url.Parse(baseHost)
		if err != nil {
			return nil, fmt.Errorf("error parsing baseHost: %w", err)
		}

		// Парсим короткий URL
		shortURL, err := url.Parse(storageURL.ShortURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing short URL: %w", err)
		}

		// Объединяем baseURL и shortURL
		fullURL := baseURL.ResolveReference(shortURL).String()

		// Сохраняем полный URL в структуру
		storageURL.ShortURL = fullURL

		urls = append(urls, storageURL)
	}

	// Проверяем на наличие ошибок после завершения обработки строк
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in row iteration: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// Возвращаем список URL-адресов
	return urls, nil
}
