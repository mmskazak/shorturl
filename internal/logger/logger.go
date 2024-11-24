package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init создает и настраивает новый логгер на основе уровня логирования.
// Принимает уровень логирования (zapcore.Level) и
// возвращает указатель на SugaredLogger и ошибку (если таковая имеется).
// Использует конфигурацию по умолчанию для производственного логгера.
func Init(level zapcore.Level) (*zap.SugaredLogger, error) {
	// Создание конфигурации для логгера с настройками по умолчанию.
	cfg := zap.NewProductionConfig()
	// Установка уровня логирования.
	cfg.Level = zap.NewAtomicLevelAt(level)

	// Построение логгера на основе конфигурации.
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("ошибка в инициализации логгера %w", err)
	}

	// Создание SugaredLogger для более удобного использования логгера с методами Sugared.
	sugar := logger.Sugar()

	return sugar, nil
}
