package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// DeleteURLs выполняет batch update записей, устанавливая флаг удаления.
//
// Функция удаляет URL-адреса, устанавливая флаг `deleted` в TRUE для каждой записи в списке `urlIDs`.
// Если список пуст, функция завершает выполнение без действий. Обновления выполняются в транзакции.
// В случае ошибки все изменения откатываются.
//
// Параметры:
// - ctx: контекст выполнения запроса.
// - urlIDs: список URL-идентификаторов для удаления.
//
// Возвращаемые значения:
// - error: ошибка, если она произошла в процессе выполнения запроса.
func (s *PostgreSQL) DeleteURLs(ctx context.Context, urlIDs []string) error {
	if len(urlIDs) == 0 {
		return nil // Если список пуст, ничего не делаем
	}

	// Начинаем транзакцию
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			s.zapLog.Errorf("Error rollback transaction: %v", err)
		}
	}()

	// Создаем batch для группового обновления
	batch := &pgx.Batch{}

	batchSize := 5000
	batchSizeCounter := 0
	// Добавляем команды в batch
	for _, shortURL := range urlIDs {
		// Если используется поле `id` в тестах, измените запрос на "WHERE id = $1"
		batch.Queue("UPDATE urls SET deleted = TRUE WHERE short_url = $1", shortURL)

		if batchSizeCounter >= batchSize {
			// Выполняем batch
			br := tx.SendBatch(ctx, batch)
			err = br.Close()
			if err != nil {
				s.zapLog.Errorf("Failed to delete URLs in batch: %v", err)
				// Откатываем транзакцию в случае ошибки
				rollbackErr := tx.Rollback(ctx)
				if rollbackErr != nil {
					s.zapLog.Errorf("Failed to rollback transaction: %v", rollbackErr)
				}
				return fmt.Errorf("failed to delete URLs in batch: %w", err)
			}
			batch = &pgx.Batch{}
			batchSizeCounter = 0
		}
		batchSizeCounter++
	}

	// Сохраняем оставшиеся данные
	if batchSizeCounter != 0 {
		// Выполняем batch
		br := tx.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			s.zapLog.Errorf("Failed to delete URLs in batch: %v", err)
			// Откатываем транзакцию в случае ошибки
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				s.zapLog.Errorf("Failed to rollback transaction: %v", rollbackErr)
			}
			return fmt.Errorf("failed to delete URLs in batch: %w", err)
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		s.zapLog.Errorf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.zapLog.Infof("Successfully deleted URLs: %v", urlIDs)
	return nil
}
