package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
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

// Сброс всех флагов, в том числе от пакета testing.
func resetFlags() {
	//nolint:reassign // переопределение флагов оправдано для тестов
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

// Инициализация конфига с флагами.
func TestInitConfigWithFlags(t *testing.T) {
	// Сбрасываем флаги перед запуском теста.
	resetFlags()

	// Устанавливаем флаги как если бы они были переданы через командную строку.
	//nolint:reassign // переопределение флагов оправдано для тестов
	os.Args = []string{"cmd", "-a", ":9090", "-b", "http://example.com", "-r", "5s", "-w", "5s", "-l", "debug",
		"-f", "/tmp/test-db.json", "-d", "postgres://user:password@localhost/db", "-secret", "newsecret"}

	// Инициализируем конфигурацию.
	config, err := InitConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем значения конфигурации.
	if config.Address != ":9090" {
		t.Errorf("Expected Address ':9090', got '%s'", config.Address)
	}
	if config.BaseHost != "http://example.com" {
		t.Errorf("Expected BaseHost 'http://example.com', got '%s'", config.BaseHost)
	}
	if config.ReadTimeout != 5*time.Second {
		t.Errorf("Expected ReadTimeout '5s', got '%s'", config.ReadTimeout)
	}
	if config.WriteTimeout != 5*time.Second {
		t.Errorf("Expected WriteTimeout '5s', got '%s'", config.WriteTimeout)
	}
	if config.LogLevel != "debug" {
		t.Errorf("Expected LogLevel 'debug', got '%s'", config.LogLevel)
	}
	if config.FileStoragePath != "/tmp/test-db.json" {
		t.Errorf("Expected FileStoragePath '/tmp/test-db.json', got '%s'", config.FileStoragePath)
	}
	if config.DataBaseDSN != "postgres://user:password@localhost/db" {
		t.Errorf("Expected DataBaseDSN 'postgres://user:password@localhost/db', got '%s'", config.DataBaseDSN)
	}
	if config.SecretKey != "newsecret" {
		t.Errorf("Expected SecretKey 'newsecret', got '%s'", config.SecretKey)
	}
}

// Тест работает только из консоли
// из-за особенностей сброса флагов
func TestInitConfigWithEnvVars(t *testing.T) {
	resetFlags() // Сбрасываем флаги командной строки.

	envVars := make(map[string]string)

	// Устанавливаем переменные окружения и очищаем их после теста.
	setEnvVars := func() {
		envVars = map[string]string{
			"SERVER_ADDRESS":    ":7070",
			"BASE_URL":          "http://env.example.com",
			"READ_TIMEOUT":      "12s",
			"WRITE_TIMEOUT":     "11s",
			"LOG_LEVEL":         "error",
			"FILE_STORAGE_PATH": "/tmp/env-db.json",
			"DATABASE_DSN":      "mysql://user:password@localhost/envdb",
			"SECRET_KEY":        "envsecret",
		}

		for key, value := range envVars {
			err := os.Setenv(key, value)
			require.NoError(t, err)
		}
	}

	defer t.Cleanup(func() {
		for key := range envVars {
			err := os.Unsetenv(key)
			require.NoError(t, err)
		}
	})

	setEnvVars()

	// Инициализируем конфигурацию.
	config, err := InitConfig()
	require.NoError(t, err)

	// Проверяем значения конфигурации.
	verifyConfig(t, config)
}

func verifyConfig(t *testing.T, config *Config) {
	t.Helper() // Указывает, что это вспомогательная функция

	if config.Address != ":7070" {
		t.Errorf("Expected Address ':7070', got '%s'", config.Address)
	}
	if config.BaseHost != "http://env.example.com" {
		t.Errorf("Expected BaseHost 'http://env.example.com', got '%s'", config.BaseHost)
	}
	if config.ReadTimeout != 12*time.Second {
		t.Errorf("Expected ReadTimeout '12s', got '%s'", config.ReadTimeout)
	}
	if config.WriteTimeout != 11*time.Second {
		t.Errorf("Expected WriteTimeout '11s', got '%s'", config.WriteTimeout)
	}
	if config.LogLevel != "error" {
		t.Errorf("Expected LogLevel 'error', got '%s'", config.LogLevel)
	}
	if config.FileStoragePath != "/tmp/env-db.json" {
		t.Errorf("Expected FileStoragePath '/tmp/env-db.json', got '%s'", config.FileStoragePath)
	}
	if config.DataBaseDSN != "mysql://user:password@localhost/envdb" {
		t.Errorf("Expected DataBaseDSN 'mysql://user:password@localhost/envdb', got '%s'", config.DataBaseDSN)
	}
	if config.SecretKey != "envsecret" {
		t.Errorf("Expected SecretKey 'envsecret', got '%s'", config.SecretKey)
	}
}
