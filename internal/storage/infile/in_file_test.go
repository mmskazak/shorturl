package infile

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mmskazak/shorturl/internal/config"

	"go.uber.org/zap"
)

func TestNewInFileErr(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	zapLog := zap.NewNop().Sugar()
	got, err := NewInFile(ctx, cfg, zapLog)
	assert.EqualError(t, err, "error read storage data: error opening or creating file: open :"+
		" no such file or directory")
	assert.Equal(t, (*InFile)(nil), got)
}

func TestNewInFileSuccessWithFile(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		Address:         ":8080",
		BaseHost:        "http://localhost:8080",
		SecretKey:       "secret",
		LogLevel:        "info",
		FileStoragePath: "storage_test.json",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
	}
	zapLog := zap.NewNop().Sugar()
	got, err := NewInFile(ctx, cfg, zapLog)
	require.NoError(t, err)
	assert.Equal(t, "storage_test.json", got.filePath)
	original, err := got.InMe.GetShortURL(ctx, "short123")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", original)
}

func Test_parseShortURLStruct_Success(t *testing.T) {
	validJSON := `{
    "id": "1",
    "short_url": "testtest",
    "original_url": "http://original.url",
    "user_id": "user123",
    "deleted": false
	}`
	sURL := shortURLStruct{
		ID:          "1",
		ShortURL:    "testtest",
		OriginalURL: "http://original.url",
		UserID:      "user123",
		Deleted:     false,
	}

	got, err := parseShortURLStruct(validJSON)
	require.NoError(t, err)
	assert.Equal(t, sURL, got)
}

func Test_parseShortURLStruct_Err(t *testing.T) {
	validJSON := `{"id": "1",`
	sURL := shortURLStruct{}

	got, err := parseShortURLStruct(validJSON)
	assert.Error(t, err)
	assert.Equal(t, sURL, got)
}
