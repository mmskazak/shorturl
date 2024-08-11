package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

// GetShortURL получает оригинальный URL по короткому URL.
//
// Функция выполняет запрос к базе данных для получения оригинального URL и флага удаления
// по короткому URL. Если короткий URL не найден, возвращает ошибку `ErrNotFound`.
// Если URL был помечен как удаленный, возвращает ошибку `ErrDeletedShortURL`.
//
// Параметры:
// - ctx: контекст выполнения запроса.
// - shortURL: короткий URL, по которому нужно получить оригинальный URL.
//
// Возвращаемые значения:
// - string: оригинальный URL, если он найден и не помечен как удаленный.
// - error: ошибка, если она произошла в процессе выполнения запроса.
func (s *PostgreSQL) GetShortURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	var deleted bool
	err := s.pool.QueryRow(ctx, "SELECT original_url, deleted FROM urls "+
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
