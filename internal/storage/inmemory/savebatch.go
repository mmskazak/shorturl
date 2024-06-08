package inmemory

import (
	"context"
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

func (m *InMemory) SaveBatch(ctx context.Context, items []storage.Incoming, baseHost string) ([]storage.Output, error) {
	dontChangedData := m.data

	outputs := make([]storage.Output, 0, len(items))
	for _, v := range items {
		err := m.SetShortURL(ctx, v.CorrelationID, v.OriginalURL)
		errUniqueViolation := errors.Is(err, storageErrors.ErrKeyAlreadyExists) ||
			errors.Is(err, storageErrors.ErrOriginalURLAlreadyExists)
		if err != nil && errUniqueViolation {
			m.data = dontChangedData
			return nil, storageErrors.ErrUniqueViolation
		}
		if err != nil {
			m.data = dontChangedData
			return nil, fmt.Errorf("save batch error: %w", err)
		}

		fullShortURL, err := storage.GetFullShortURL(baseHost, v.CorrelationID)
		if err != nil {
			m.data = dontChangedData
			return nil, fmt.Errorf("error getFullShortURL from two parts %w", err)
		}

		outputs = append(outputs, storage.Output{
			CorrelationID: v.CorrelationID,
			ShortURL:      fullShortURL,
		})
	}

	return outputs, nil
}
