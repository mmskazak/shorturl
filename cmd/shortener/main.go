package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"mmskazak/shorturl/internal/services/shorturlservice"

	"mmskazak/shorturl/internal/app"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/storage/factory"
)

// main инициализирует конфигурацию, логгер, хранилище и запускает приложение.
//
//go:generate go run ./../version/main.go
func main() {
	// Инициализация конфигурации.
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Ошибка инициализации конфигурации: %v", err)
	}

	// Получение уровня логирования из конфигурации.
	level, err := cfg.LogLevel.Value()
	if err != nil {
		log.Printf("Ошибка получения уровня логирования: %v", err)
	}

	// Инициализация логгера.
	zapLog, err := logger.Init(level)
	if err != nil {
		log.Printf("ошибка инициализации логера output: %v", err)
	}

	// Создание контекста.
	ctx := context.Background()

	// Инициализация хранилища.
	storage, err := factory.NewStorage(ctx, cfg, zapLog)
	if err != nil {
		zapLog.Fatalf("Ошибка инициализации хранилища: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			zapLog.Warn("Error closing storage: %v\n", err)
		}
	}()

	shortURLService := shorturlservice.NewShortURLService()

	// Создание и запуск приложения.
	newApp := app.NewApp(
		ctx,
		cfg,
		storage,
		cfg.ReadTimeout,
		cfg.WriteTimeout,
		zapLog,
		shortURLService,
	)

	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zapLog.Fatalf("Ошибка сервера: %v", err)
	}

	zapLog.Infoln("Приложение завершило работу.")
}
