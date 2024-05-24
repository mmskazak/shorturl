package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	PermFile0644 = 0o644
)

var Logf *zap.SugaredLogger

func InitWriteToFile(level zapcore.Level) (*zap.Logger, error) {
	// Open files for logging
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, PermFile0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	errorFile, err := os.OpenFile("app_error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, PermFile0644)
	if err != nil {
		err := file.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close simple log file: %w", err)
		}
		return nil, fmt.Errorf("failed to open error log file: %w", err)
	}

	// Configure logger
	cfg := zap.Config{
		Encoding:         "json",                           // Log format (json, console, etc.)
		Level:            zap.NewAtomicLevelAt(level),      // Logging level
		OutputPaths:      []string{"app.log"},              // Log file path
		ErrorOutputPaths: []string{"app_error.log"},        // Error log file path
		EncoderConfig:    zap.NewProductionEncoderConfig(), // Encoder configuration
	}

	// Build logger with configuration
	logger, err := cfg.Build()
	if err != nil {
		err := file.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close simple log file: %w", err)
		}
		err = errorFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close error log file: %w", err)
		}
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	// Ensure resources are closed properly
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("failed to close log file: %w", err)
	}

	if err := errorFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close error log file: %w", err)
	}

	// Ensure logger sync
	if err := logger.Sync(); err != nil {
		return nil, fmt.Errorf("failed to sync logger: %w", err)
	}

	Logf = logger.Sugar()

	return logger, nil
}
