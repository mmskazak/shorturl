package main

import (
	"context"
	"errors"
	"log"
	"mmskazak/shorturl/internal/services/shorturlservice"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ctx, cancel := context.WithCancel(ctx)

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

	zapLog.Infof("Build version: %s", buildVersion)
	zapLog.Infof("Build date: %s", buildDate)
	zapLog.Infof("Build commit: %s", buildCommit)

	// Создаем канал для получения системных сигналов.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Запуск сервера в отдельной горутине.
	go func() {
		if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zapLog.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидаем сигнал завершения.
	<-quit
	zapLog.Infoln("Получен сигнал завершения, остановка сервера...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := newApp.Stop(ctxShutdown); err != nil {
		zapLog.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	zapLog.Infoln("Приложение завершило работу.")
}
