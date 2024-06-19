package infile

import (
	"context"
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
	m.saveToFile()
	m.zapLog.Infof("Added short link id: %s, URL: %s, for UserID %s", idShortPath, originalURL, userID)
	return nil
}
