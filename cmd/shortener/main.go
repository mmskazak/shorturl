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
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("ошибка инициализации конфигурации в main %v", err)
	}

	level, err := cfg.LogLevel.Value()
	if err != nil {
		log.Printf("ошибка получения уровня логирония %v", err)
	}
	_, err = logger.Init(level)
	if err != nil {
		log.Fatalf("ошибка инициализации логера в main %v", err)
	}

	var ms *mapstorage.MapStorage
	pathToStorage := cfg.FileStoragePath
	if pathToStorage != "" {
		logger.Log.Debug(pathToStorage)
		// filepath.Split возвращает пару (директория, имя файла), но нам нужно только директорию
		directoryPath, _ := filepath.Split(pathToStorage)
		err = os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			log.Fatalf("ошибка создания дирректории для файла хранилища %v", err)
		}

		consumer, err := rwstorage.NewConsumer(pathToStorage)
		if err != nil {
			log.Fatalf("ошибка создания консьюмера для чтения из хранилища")
		}

		ms = mapstorage.NewMapStorage(cfg.FileStoragePath)

		// Читаем события в цикле
		for {
			dataOfURL, err := consumer.ReadDataFromFile()
			if err != nil {
				if err.Error() != "EOF" {
					fmt.Printf("Ошибка при чтении: %v\n", err)
					break
				}
				fmt.Println("Достигнут конец файла.")
				break
			}

			// Обрабатываем прочитанное событие
			fmt.Printf("Прочитанны данные: %+v\n", dataOfURL)
			ms.Data[dataOfURL.ShortURL] = dataOfURL.OriginalURL
			fmt.Printf("длинна мапы: %+v\n", len(ms.Data))
		}
	} else {
		ms = mapstorage.NewMapStorage(cfg.FileStoragePath)
	}

	newApp := app.NewApp(cfg, ms, cfg.ReadTimeout, cfg.WriteTimeout)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}

	log.Println("Application has shut down gracefully.")
}
