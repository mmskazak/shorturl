package infile

import (
	"context"
	"mmskazak/shorturl/internal/models"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/storage/inmemory"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInFile_SaveBatch(t *testing.T) {
	ctx := context.Background()
	// Создание экземпляра InMemory
	inMemory, err := inmemory.NewInMemory(zap.NewNop().Sugar())
	require.NoError(t, err)
	incoming := []models.Incoming{
		{
			CorrelationID: "123",
			OriginalURL:   "https://example.com/long-url-00012",
		},
		{
			CorrelationID: "456",
			OriginalURL:   "https://example.com/long-url-00013",
		},
	}

	// Базовый хост
	baseHost := "http://localhost"
	userID := "1"
	generator := genidurl.NewGenIDService()

	s := &InFile{
		InMe:     inMemory,
		zapLog:   zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
		filePath: "/mock/path",
	}
	got, err := s.SaveBatch(ctx, incoming, baseHost, userID, generator)
	require.NoError(t, err)
	assert.Equal(t, 2, len(got))
}
