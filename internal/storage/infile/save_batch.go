package infile

import (
	"context"

	"mmskazak/shorturl/internal/storage"
)

func (m *InFile) SaveBatch(
	ctx context.Context,
	items []storage.Incoming,
	baseHost string,
	userID string,
	generator storage.IGenIDForURL,
) ([]storage.Output, error) {
	outputs, err := m.InMe.SaveBatch(ctx, items, baseHost, userID, generator)
	if err != nil {
		return nil, err //nolint:wrapcheck // прокидываем оригинальную ошибку
	}
	m.saveToFile()
	return outputs, nil
}
