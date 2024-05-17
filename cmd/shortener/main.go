package main

import (
	"errors"
	"log"
	"mmskazak/shorturl/internal/app"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/logger"
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

	newApp := app.NewApp(cfg, ms, cfg.ReadTimeout, cfg.WriteTimeout)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
