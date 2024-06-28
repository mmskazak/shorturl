package postgresql

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

var timeSleep = 24 * time.Hour

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
