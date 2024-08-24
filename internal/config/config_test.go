package config

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestLogLevel_Value(t *testing.T) {
	tests := []struct {
		name          string
		logLevel      LogLevel
		expectedLevel zapcore.Level
		expectError   bool
	}{
		{"Debug level", "debug", zapcore.DebugLevel, false},
		{"Info level", "info", zapcore.InfoLevel, false},
		{"Warn level", "warn", zapcore.WarnLevel, false},
		{"Warning level", "warning", zapcore.WarnLevel, false},
		{"Error level", "error", zapcore.ErrorLevel, false},
		{"DPanic level", "dpanic", zapcore.DPanicLevel, false},
		{"Panic level", "panic", zapcore.PanicLevel, false},
		{"Fatal level", "fatal", zapcore.FatalLevel, false},
		{"Unknown level", "unknown", zapcore.DebugLevel, true},
		{"Empty level", "", zapcore.DebugLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := tt.logLevel.Value()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedLevel, level)
		})
	}
}
