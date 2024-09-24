package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
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
