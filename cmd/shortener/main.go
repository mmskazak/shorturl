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
	"strconv"
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

	pathToStorage := cfg.FileStoragePath

	_, err = logger.Init(level)
	if err != nil {
		log.Fatalf("ошибка инициализации логера в main %v", err)
	}

	// filepath.Split возвращает пару (директория, имя файла), но нам нужно только директорию
	directoryPath, filename := filepath.Split(pathToStorage)
	err = os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
		log.Fatalf("ошибка создания дирректории для файла хранилища %v", err)
	}

	consumer, err := rwstorage.NewConsumer(filename)
	if err != nil {
		log.Fatalf("ошибка создания консьюмера для чтения из хранилища")
	}

	ms := mapstorage.NewMapStorage()

	// Читаем события в цикле
	for {
		dataOfURL, err := consumer.ReadDataFromFile()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("Ошибка при чтении события: %v\n", err)
				break
			}
			fmt.Println("Достигнут конец файла.")
			break
		}

		// Обрабатываем прочитанное событие
		fmt.Printf("Прочитанны данные: %+v\n", dataOfURL)
		err = ms.SetShortURL(dataOfURL.ShortURL, dataOfURL.OriginalURL)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	newApp := app.NewApp(cfg, ms, cfg.ReadTimeout, cfg.WriteTimeout)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}

	defer recordMapToFile(ms, filename)
	log.Println("Application has shut down gracefully.")
}

func recordMapToFile(ms *mapstorage.MapStorage, filename string) {
	producer, err := rwstorage.NewProducer(filename)
	if err != nil {
		log.Fatalf("ошибка создания продюсера %v", err)
	}
	counter := 0
	for value, key := range ms.Data {
		counter++
		shul := rwstorage.ShortURLStruct{
			UUID:        strconv.Itoa(counter),
			ShortURL:    key,
			OriginalURL: value,
		}

		err := producer.WriteData(&shul)
		if err != nil {
			log.Fatal(err)
		}
	}
	producer.Close()
	log.Println("Cleanup completed.")
}
