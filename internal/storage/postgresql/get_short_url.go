package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

func (s *PostgreSQL) GetShortURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	var deleted bool
	err := s.pool.QueryRow(ctx, "SELECT original_url FROM urls "+
		"WHERE short_url = $1", shortURL).Scan(&originalURL, &deleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("short URL not found: %w", storageErrors.ErrNotFound)
		}
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}
	if deleted {
		return "", storageErrors.ErrDeletedShortURL
	}
	return originalURL, nil
}
