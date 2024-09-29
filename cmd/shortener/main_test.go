package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage/inmemory"
	"testing"
	"time"
)

func Test_prepareParamsForApp(t *testing.T) {
	ctx := context.Background()
	cfg, zapLog, storage := prepareParamsForApp(ctx)
	assert.Equal(t, "/tmp/short-url-db.json", cfg.FileStoragePath)
	assert.NotNil(t, zapLog)
	assert.NotNil(t, storage)
}

func Test_loggingBuildParams(t *testing.T) {
	zapSugar := zaptest.NewLogger(t).Sugar()
	loggingBuildParams(zapSugar)
}

func TestRunApp(t *testing.T) {
	ctx := context.Background()
	ctxWt, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	baseDurationReadTimeout := 10 * time.Second  //nolint:gomnd  // 10 секунд.
	baseDurationWriteTimeout := 10 * time.Second //nolint:gomnd  // 10 секунд.
	cfg := &config.Config{
		Address:         ":8080",
		BaseHost:        "http://localhost:8080",
		LogLevel:        "info",
		ReadTimeout:     baseDurationReadTimeout,
		WriteTimeout:    baseDurationWriteTimeout,
		FileStoragePath: "/tmp/short-url-db.json",
		SecretKey:       "secret",
		ConfigPath:      "",
		TrustedSubnet:   "",
	}

	zapLog := zap.NewNop().Sugar()
	storage, err := inmemory.NewInMemory(zapLog)

	shutdownDuration := 5 * time.Second

	err = runApp(ctxWt, cfg, zapLog, storage, shutdownDuration)
	require.NoError(t, err)
}
