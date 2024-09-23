package inmemory

import (
	"context"
	"mmskazak/shorturl/internal/models"
	"strconv"
)

// InternalStats - count users and urls in inmemory storage
func (m *InMemory) InternalStats(_ context.Context) (models.Stats, error) {
	m.zapLog.Info("Getting internal stats from InMemory store.")
	countUrls := len(m.data)
	countUsers := len(m.userIndex)

	stats := models.Stats{
		Urls:  strconv.Itoa(countUrls),
		Users: strconv.Itoa(countUsers),
	}

	return stats, nil
}
