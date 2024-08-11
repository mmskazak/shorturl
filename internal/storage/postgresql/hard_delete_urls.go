package postgresql

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
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
// - db: пул соединений с базой данных pgxpool.Pool.
// - zapLog: логгер для записи ошибок и информационных сообщений.
func hardDeleteSoftDeletedURLs(ctx context.Context, db *pgxpool.Pool, zapLog *zap.SugaredLogger) {
	for {
		query := `DELETE FROM urls WHERE deleted=true`
		_, err := db.Exec(ctx, query)
		if err != nil {
			zapLog.Errorf("error deleting urls: %v", err)
		}

		time.Sleep(timeSleep)
	}
}
