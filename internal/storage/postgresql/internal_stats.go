package postgresql

import (
	"context"
	"fmt"

	"mmskazak/shorturl/internal/models"

	"github.com/jackc/pgx/v5"
)

// InternalStats - count users and urls in a database.
func (s *PostgreSQL) InternalStats(ctx context.Context) (models.Stats, error) {
	queryCountUrls := "SELECT COUNT(original_url) FORM urls WHERE 1"
	queryCountUsers := "SELECT COUNT(user_id) FORM urls WHERE 1 ORDER BY (user_id)"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.Stats{}, fmt.Errorf("error starting transaction for internal stats %w", err)
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			s.zapLog.Errorf("Error rollback transaction: %v", err)
		}
	}(tx, ctx)

	countUrls := tx.QueryRow(ctx, queryCountUrls)
	var urls string
	err = countUrls.Scan(&urls)
	if err != nil {
		return models.Stats{}, fmt.Errorf("error query row for count urls %w", err)
	}

	row := tx.QueryRow(ctx, queryCountUsers)
	var users string
	err = row.Scan(row, &users)
	if err != nil {
		return models.Stats{}, fmt.Errorf("error query row for count users %w", err)
	}

	stats := models.Stats{
		Urls:  urls,
		Users: users,
	}

	return stats, nil
}
