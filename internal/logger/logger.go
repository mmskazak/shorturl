package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(level zapcore.Level) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build()

	if err != nil {
		return nil, fmt.Errorf("ошибка в инициализации логера %w", err)
	}

	sugar := logger.Sugar()

	return sugar, nil
}
