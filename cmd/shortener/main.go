package main

import (
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/app"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/services/rwstorage"
	"mmskazak/shorturl/internal/storage/mapstorage"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const (
	PermFile0750 = 0o750
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Ошибка инициализации конфигурации: %v", err)
	}

	level, err := cfg.LogLevel.Value()
	if err != nil {
		log.Printf("Ошибка получения уровня логирования: %v", err)
	}

	zapLog, err := logger.Init(level)
	if err != nil {
		log.Printf("ошибка инициализации логера output: %v", err)
	}

	ms, err := initializeStorage(cfg, zapLog)
	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	newApp := app.NewApp(cfg, ms, cfg.ReadTimeout, cfg.WriteTimeout, zapLog)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Ошибка сервера: %v", err)
	}

	log.Println("Приложение завершило работу.")
}

func initializeStorage(cfg *config.Config, zapLog *zap.SugaredLogger) (*mapstorage.MapStorage, error) {
	pathToStorage := cfg.FileStoragePath
	var ms *mapstorage.MapStorage

	if pathToStorage == "" {
		return mapstorage.NewMapStorage(pathToStorage), nil
	}

	if err := os.MkdirAll(filepath.Dir(pathToStorage), os.FileMode(PermFile0750)); err != nil {
		return nil, fmt.Errorf("ошибка создания директории для файла хранилища: %w", err)
	}

	consumer, err := rwstorage.NewConsumer(pathToStorage)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания консьюмера для чтения из хранилища: %w", err)
	}

	ms = mapstorage.NewMapStorage(pathToStorage)
	if err := readFileStorage(consumer, ms, zapLog); err != nil {
		return nil, fmt.Errorf("error read storage data: %w", err)
	}

	return ms, nil
}

func readFileStorage(consumer *rwstorage.Consumer, ms *mapstorage.MapStorage, zapLog *zap.SugaredLogger) error {
	for {
		dataOfURL, err := consumer.ReadDataFromFile()
		if err != nil {
			if err.Error() != "EOF" {
				return fmt.Errorf("ошибка при чтении: %w", err)
			}
			fmt.Println("Достигнут конец файла.")
			break
		}

		zapLog.Infof("Прочитанные данные: %+v\n", dataOfURL)
		ms.Data[dataOfURL.ShortURL] = dataOfURL.OriginalURL
		zapLog.Infof("Длина мапы: %+v\n", len(ms.Data))
	}
	return nil
}
