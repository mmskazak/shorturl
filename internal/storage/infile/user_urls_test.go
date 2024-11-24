package infile

import (
	"context"
	"testing"

	"mmskazak/shorturl/internal/models"

	"mmskazak/shorturl/internal/storage/inmemory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInFile_GetUserURLs_Err(t *testing.T) {
	ctx := context.Background()
	userID := "1"
	baseHost := ""
	inm, err := inmemory.NewInMemory(zap.NewNop().Sugar())
	require.NoError(t, err)
	err = inm.SetShortURL(context.Background(),
		"short123",
		"https://example.com",
		"1",
		false)
	require.NoError(t, err)
	m := &InFile{
		InMe:     inm,
		zapLog:   zap.NewNop().Sugar(),
		filePath: "/mock/path",
	}
	got, err := m.GetUserURLs(ctx, userID, baseHost)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestInFile_GetUserURLs_Success(t *testing.T) {
	ctx := context.Background()
	userID := "1"
	baseHost := "http://localhost"
	inm, err := inmemory.NewInMemory(zap.NewNop().Sugar())
	require.NoError(t, err)
	err = inm.SetShortURL(context.Background(),
		"short123",
		"https://example.com",
		"1",
		false)
	require.NoError(t, err)
	m := &InFile{
		InMe:     inm,
		zapLog:   zap.NewNop().Sugar(),
		filePath: "/mock/path",
	}
	got, err := m.GetUserURLs(ctx, userID, baseHost)
	assert.NoError(t, err)
	assert.Equal(t, []models.URL{
		{
			ShortURL:    "http://localhost/short123",
			OriginalURL: "https://example.com"},
	}, got)
}
