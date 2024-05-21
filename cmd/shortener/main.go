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

	ms := mapstorage.NewMapStorage()

	filename := "short-url-db.json"
	// Чтение содержимого файла
	consumer, _ := rwstorage.NewConsumer(filename)

	// Читаем события в цикле
	for {
		dataOfURL, err := consumer.ReadEvent()
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
}
