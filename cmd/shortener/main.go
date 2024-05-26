package main

import (
	"errors"
	"log"
	"mmskazak/shorturl/internal/app"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/storage"
	"net/http"
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

	var storageType string
	switch {
	case cfg.FileStoragePath == "":
		storageType = "inmemory"
	case cfg.FileStoragePath != "":
		storageType = "infile"
	}

	ms, err := storage.NewStorage(storageType, cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	newApp := app.NewApp(cfg, ms, cfg.ReadTimeout, cfg.WriteTimeout, zapLog)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Ошибка сервера: %v", err)
	}

	log.Println("Приложение завершило работу.")
}
