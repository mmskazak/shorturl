package main

import (
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/app"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/services/rwstorage"
	"mmskazak/shorturl/internal/storage/mapstorage"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	PermFile0750 = 0o750
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Logf.Fatalf("Ошибка инициализации конфигурации: %v", err)
	}

	level, err := cfg.LogLevel.Value()
	if err != nil {
		logger.Logf.Warnf("Ошибка получения уровня логирования: %v", err)
	}

	if err := initLoggers(level); err != nil {
		logger.Logf.Fatalf("Ошибка инициализации логеров: %v", err)
	}

	ms, err := initializeStorage(cfg, logger.Logf)
	if err != nil {
		logger.Logf.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	newApp := app.NewApp(cfg, ms, cfg.ReadTimeout, cfg.WriteTimeout)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Logf.Fatalf("Ошибка сервера: %v", err)
	}

	logger.Logf.Info("Приложение завершило работу.")
}

func initLoggers(level zapcore.Level) error {
	if _, err := logger.InitWriteToOutput(level); err != nil {
		return fmt.Errorf("ошибка инициализации логера output: %w", err)
	}
	if _, err := logger.InitWriteToFile(level); err != nil {
		return fmt.Errorf("ошибка инициализации логера file: %w", err)
	}
	return nil
}

func initializeStorage(cfg *config.Config, logf *zap.SugaredLogger) (*mapstorage.MapStorage, error) {
	pathToStorage := cfg.FileStoragePath
	var ms *mapstorage.MapStorage

	if pathToStorage == "" {
		return mapstorage.NewMapStorage(pathToStorage), nil
	}

	logf.Debug(pathToStorage)
	if err := os.MkdirAll(filepath.Dir(pathToStorage), os.FileMode(PermFile0750)); err != nil {
		return nil, fmt.Errorf("ошибка создания директории для файла хранилища: %w", err)
	}

	consumer, err := rwstorage.NewConsumer(pathToStorage)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания консьюмера для чтения из хранилища: %w", err)
	}

	ms = mapstorage.NewMapStorage(pathToStorage)
	if err := readFileStorage(consumer, ms, logf); err != nil {
		return nil, fmt.Errorf("error read storage data: %w", err)
	}

	return ms, nil
}

func readFileStorage(consumer *rwstorage.Consumer, ms *mapstorage.MapStorage, logf *zap.SugaredLogger) error {
	for {
		dataOfURL, err := consumer.ReadDataFromFile()
		if err != nil {
			if err.Error() != "EOF" {
				return fmt.Errorf("ошибка при чтении: %w", err)
			}
			fmt.Println("Достигнут конец файла.")
			break
		}

		logf.Infof("Прочитанные данные: %+v\n", dataOfURL)
		ms.Data[dataOfURL.ShortURL] = dataOfURL.OriginalURL
		logf.Infof("Длина мапы: %+v\n", len(ms.Data))
	}
	return nil
}
