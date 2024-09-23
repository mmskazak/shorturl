package infile

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/models"
)

// InternalStats - count users and urls in inmemory storage
func (m *InFile) InternalStats(ctx context.Context) (models.Stats, error) {
	stats, err := m.InternalStats(ctx)
	if err != nil {
		return models.Stats{}, fmt.Errorf("error getting internal stats form file %w", err)
	}

	return stats, nil
}
