package main

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"log"
	"mmskazak/shorturl/internal/contracts"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mmskazak/shorturl/internal/services/shorturlservice"

	"mmskazak/shorturl/internal/app"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/storage/factory"
)

const shutdownDuration = 5 * time.Second

// main инициализирует конфигурацию, логгер, хранилище и запускает приложение.
//
//go:generate go run ./../version/main.go
func main() {
	// Создание контекста.
	ctx := context.Background()

	cfg, zapLog, storage := prepareParamsForApp(ctx)

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

	// Логирование Build параметров
	loggingBuildParams(zapLog)

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

	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownDuration)
	defer cancel()

	if err := newApp.Stop(ctxShutdown); err != nil {
		zapLog.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	zapLog.Infoln("Приложение завершило работу.")
}

func prepareParamsForApp(ctx context.Context) (*config.Config, *zap.SugaredLogger, contracts.Storage) {
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

	// Инициализация хранилища.
	storage, err := factory.NewStorage(ctx, cfg, zapLog)
	if err != nil {
		zapLog.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	return cfg, zapLog, storage
}

func loggingBuildParams(zapLog *zap.SugaredLogger) {
	zapLog.Infof("Build version: %s", buildVersion)
	zapLog.Infof("Build date: %s", buildDate)
	zapLog.Infof("Build commit: %s", buildCommit)
}
