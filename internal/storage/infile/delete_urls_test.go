package infile

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/storage/inmemory"
	"testing"
)

func TestInFile_DeleteURLs(t *testing.T) {
	ctx := context.Background()
	urlIDs := []string{
		"testtest",
		"test1234",
		"1234test",
	}
	// Создание экземпляра InMemory
	inMemory, err := inmemory.NewInMemory(zap.NewNop().Sugar())
	require.NoError(t, err)

	// Создание экземпляра InFile
	s := &InFile{
		InMe:     inMemory,
		zapLog:   zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
		filePath: "/mock/path",
	}

	err = s.DeleteURLs(ctx, urlIDs)
	require.NoError(t, err)
}
