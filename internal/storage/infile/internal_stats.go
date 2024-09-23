package infile

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/models"
)

// InternalStats - count users and urls in inmemory storage
func (f *InFile) InternalStats(ctx context.Context) (models.Stats, error) {
	f.zapLog.Info("Getting internal stats from InFile store.")
	stats, err := f.InMe.InternalStats(ctx)
	if err != nil {
		return models.Stats{}, fmt.Errorf("error getting internal stats form file %w", err)
	}

	return stats, nil
}
