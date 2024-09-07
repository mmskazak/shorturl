package inmemory

import (
	"context"
	"errors"
	"fmt"

	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/dtos"
	"mmskazak/shorturl/internal/models"

	"mmskazak/shorturl/internal/services/shorturlservice"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

// IGenIDForURL - интерфейс генерирования ID ссылки для URL.
type IGenIDForURL interface {
	Generate() (string, error)
}

// SaveBatch error:
// different error
// ErrKeyAlreadyExists
// ConflictError (ErrOriginalURLAlreadyExists).
func (m *InMemory) SaveBatch(
	ctx context.Context,
	items []models.Incoming,
	baseHost string,
	userID string,
	generator contracts.IGenIDForURL,
) ([]models.Output, error) {
	dontChangedData := m.data

	outputs := make([]models.Output, 0, len(items))
	for _, v := range items {
		dto := dtos.DTOShortURL{
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
			m.data = dontChangedData
			return nil, conflictErr
		}
		if errors.Is(err, storageErrors.ErrKeyAlreadyExists) {
			m.data = dontChangedData
			return nil, storageErrors.ErrKeyAlreadyExists
		}
		if err != nil {
			return nil, fmt.Errorf("error inserting data: %w", err)
		}

		outputs = append(outputs, models.Output{
			CorrelationID: v.CorrelationID,
			ShortURL:      fullShortURL,
		})
	}

	return outputs, nil
}
