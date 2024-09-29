package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mmskazak/shorturl/internal/contracts"

	"go.uber.org/zap"

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
	ctx := context.Background()
	cfg, zapLog, storage := prepareParamsForApp(ctx)

	if err := runApp(ctx, cfg, zapLog, storage, shutdownDuration); err != nil {
		zapLog.Fatal("Ошибка приложения: %v", err)
	}
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

func runApp(
	ctx context.Context,
	cfg *config.Config,
	zapLog *zap.SugaredLogger,
	store contracts.Storage,
	shutdownDuration time.Duration,
) error {
	defer func() {
		if err := store.Close(); err != nil {
			zapLog.Error("Error closing store\n", zap.Error(err))
		}
	}()

	shortURLService := shorturlservice.NewShortURLService()

	newApp := app.NewApp(
		ctx,
		cfg,
		store,
		cfg.ReadTimeout,
		cfg.WriteTimeout,
		zapLog,
		shortURLService,
	)

	loggingBuildParams(zapLog)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := newApp.StartAll(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zapLog.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание завершения
	select {
	case <-quit: // Ожидание сигнала завершения
		zapLog.Infoln("Получен сигнал завершения, остановка сервера...")
	case <-ctx.Done(): // Завершение по контексту
		zapLog.Infoln("Контекст завершён, остановка сервера...")
	}

	zapLog.Infoln("Получен сигнал завершения, остановка сервера...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownDuration)
	defer cancel()

	if err := newApp.Stop(ctxShutdown); err != nil {
		zapLog.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	zapLog.Infoln("Приложение завершило работу.")
	return nil
}
