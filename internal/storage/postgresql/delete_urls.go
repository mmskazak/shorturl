package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

// DeleteURLs выполняет batch update записей, устанавливая флаг удаления.
func (s *PostgreSQL) DeleteURLs(urlIDs []string) error {
	if len(urlIDs) == 0 {
		return nil // Если список пуст, ничего не делаем
	}

	// Создаем batch для группового обновления
	batch := &pgx.Batch{}

	// Добавляем команды в batch
	for _, id := range urlIDs {
		batch.Queue("UPDATE urls SET deleted = TRUE WHERE id = $1", id)
	}

	// Выполняем batch
	br := s.pool.SendBatch(context.Background(), batch)
	err := br.Close()
	if err != nil {
		log.Printf("Failed to delete URLs in batch: %v", err)
		return err
	}

	log.Printf("Successfully deleted URLs: %v", urlIDs)
	return nil
}
