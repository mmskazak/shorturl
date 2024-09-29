package infile

import (
	"context"
	"mmskazak/shorturl/internal/models"
	"mmskazak/shorturl/internal/storage/inmemory"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInFile_InternalStatsZeroZero(t *testing.T) {
	ctx := context.Background()
	zapLog := zap.NewNop().Sugar()
	filePath := "/tmp/testpath.json"
	// Создание экземпляра InMemory
	inMemory, err := inmemory.NewInMemory(zap.NewNop().Sugar())
	if err != nil {
		t.Fatalf("Failed to create InMemory instance: %v", err)
	}

	f := &InFile{
		InMe:     inMemory,
		zapLog:   zapLog,
		filePath: filePath,
	}
	got, err := f.InternalStats(ctx)
	expected := models.Stats{Urls: "0", Users: "0"}
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestInFile_InternalStatsTreeTwo(t *testing.T) {
	ctx := context.Background()
	zapLog := zap.NewNop().Sugar()
	filePath := "/tmp/testpath.json"
	// Создание экземпляра InMemory
	inMemory, err := inmemory.NewInMemory(zap.NewNop().Sugar())
	require.NoError(t, err)
	err = inMemory.SetShortURL(ctx, "TestYanD", "http://ya.ru", "1", false)
	require.NoError(t, err)
	err = inMemory.SetShortURL(ctx, "TeStYanD", "http://yandex.ru", "1", false)
	require.NoError(t, err)
	err = inMemory.SetShortURL(ctx, "GooGLeeE", "http://google.com", "2", false)
	require.NoError(t, err)

	f := &InFile{
		InMe:     inMemory,
		zapLog:   zapLog,
		filePath: filePath,
	}
	got, err := f.InternalStats(ctx)
	expected := models.Stats{Urls: "3", Users: "2"}
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}
