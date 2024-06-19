package infile

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/services/rwstorage"
	"strconv"
)

// SetShortURL error:
// different error
// ErrKeyAlreadyExists
// ConflictError (ErrOriginalURLAlreadyExists).
func (m *InFile) SetShortURL(ctx context.Context, idShortPath string, originalURL string, userID string) error {
	err := m.InMe.SetShortURL(ctx, idShortPath, originalURL, userID)
	if err != nil {
		return err //nolint:wrapcheck // пробрасываем дальше оригиральную ошибку
	}

	record := rwstorage.ShortURLStruct{
		ID:          strconv.Itoa(m.InMe.NumberOfEntries()),
		ShortURL:    idShortPath,
		OriginalURL: originalURL,
		UserID:      userID,
		Deleted:     false,
	}

	if err := m.appendToFile(record); err != nil {
		return fmt.Errorf("error appending to file: %w", err)
	}

	m.zapLog.Infof("Added short link id: %s, URL: %s, for UserID %s", idShortPath, originalURL, userID)
	return nil
}
