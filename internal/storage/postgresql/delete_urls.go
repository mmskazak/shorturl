package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// DeleteURLs выполняет batch update записей, устанавливая флаг удаления.
func (s *PostgreSQL) DeleteURLs(ctx context.Context, urlIDs []string) error {
	if len(urlIDs) == 0 {
		return nil // Если список пуст, ничего не делаем
	}

	// Начинаем транзакцию
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Создаем batch для группового обновления
	batch := &pgx.Batch{}

	// Добавляем команды в batch
	for _, shortURL := range urlIDs {
		// Если используется поле `id` в тестах, измените запрос на "WHERE id = $1"
		batch.Queue("UPDATE urls SET deleted = TRUE WHERE short_url = $1", shortURL)
	}

	// Выполняем batch в контексте транзакции
	br := tx.SendBatch(ctx, batch)
	err = br.Close()
	if err != nil {
		log.Printf("Failed to delete URLs in batch: %v", err)
		// Откатываем транзакцию в случае ошибки
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			log.Printf("Failed to rollback transaction: %v", rollbackErr)
		}
		return fmt.Errorf("failed to delete URLs in batch: %w", err)
	}

	// Фиксируем транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully deleted URLs: %v", urlIDs)
	return nil
}
