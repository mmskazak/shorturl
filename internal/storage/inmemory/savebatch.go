package inmemory

import (
	"fmt"
	"mmskazak/shorturl/internal/storage"
)

func (m *InMemory) SaveBatch(items []storage.Incoming, baseHost string) ([]storage.Output, error) {
	dontChangedData := m.data

	outputs := make([]storage.Output, 0, len(items))
	for _, v := range items {
		err := m.SetShortURL(v.CorrelationID, v.OriginalURL)
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
