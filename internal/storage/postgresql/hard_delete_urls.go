package postgresql

import (
	"context"
	"time"

	"mmskazak/shorturl/internal/storage/postgresql/interfaces"

	"go.uber.org/zap"
)

// timeSleep определяет интервал времени между запусками удаления URL.
var timeSleep = 24 * time.Hour

// hardDeleteSoftDeletedURLs выполняет постоянное удаление записей с флагом удаления из базы данных.
//
// Функция запускает бесконечный цикл, в котором выполняется SQL-запрос на удаление записей,
// помеченных как удаленные (deleted=true). Если запрос завершился ошибкой, она логируется.
// После каждого выполнения запроса функция ожидает интервал времени, определенный переменной timeSleep,
// и повторяет процесс.
//
// Параметры:
// - ctx: контекст выполнения запроса.
// - db: пул соединений с базой данных db interfaces.Database.
// - zapLog: логгер для записи ошибок и информационных сообщений.
func hardDeleteSoftDeletedURLs(ctx context.Context, db interfaces.Database, zapLog *zap.SugaredLogger) {
	for {
		select {
		case <-ctx.Done():
			// Завершаем работу при отмене контекста
			zapLog.Info("Context done, stopping hard delete operation.")
			return
		default:
			query := `DELETE FROM urls WHERE deleted = true`
			_, err := db.Exec(ctx, query)
			if err != nil {
				zapLog.Errorf("Error deleting URLs: %v", err)
				// Возможно, стоит добавить логику для обработки ошибок или увеличить время ожидания
			} else {
				zapLog.Info("Successfully deleted soft-deleted URLs.")
			}

			// Пауза между попытками
			time.Sleep(timeSleep)
		}
	}
}
