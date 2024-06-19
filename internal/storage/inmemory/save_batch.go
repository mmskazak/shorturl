package inmemory

import (
	"context"
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/services/shorturlservice"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

type IGenIDForURL interface {
	Generate() (string, error)
}

// SaveBatch error:
// different error
// ErrKeyAlreadyExists
// ConflictError (ErrOriginalURLAlreadyExists).
func (m *InMemory) SaveBatch(
	ctx context.Context,
	items []storage.Incoming,
	baseHost string,
	userID string,
	generator storage.IGenIDForURL,
) ([]storage.Output, error) {
	dontChangedData := m.Data

	outputs := make([]storage.Output, 0, len(items))
	for _, v := range items {
		dto := shorturlservice.DTOShortURL{
			OriginalURL: v.OriginalURL,
			UserID:      userID,
			BaseHost:    baseHost,
			Deleted:     false,
		}
		service := shorturlservice.NewShortURLService()
		fullShortURL, err := service.GenerateShortURL(
			ctx,
			dto,
			generator,
			m,
		)

		var conflictErr storageErrors.ConflictError
		if errors.As(err, &conflictErr) {
			m.Data = dontChangedData
			return nil, conflictErr
		}
		if errors.Is(err, storageErrors.ErrKeyAlreadyExists) {
			m.Data = dontChangedData
			return nil, storageErrors.ErrKeyAlreadyExists
		}
		if err != nil {
			return nil, fmt.Errorf("error inserting Data: %w", err)
		}

		outputs = append(outputs, storage.Output{
			CorrelationID: v.CorrelationID,
			ShortURL:      fullShortURL,
		})
	}

	return outputs, nil
}
