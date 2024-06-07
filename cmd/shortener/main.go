package main

import (
	"context"
	"errors"
	"log"
	"mmskazak/shorturl/internal/app"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/storage/factory"
	"net/http"
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

	ctx := context.Background()
	storage, err := factory.NewStorage(ctx, cfg)
	if err != nil {
		zapLog.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	newApp := app.NewApp(ctx, cfg, storage, cfg.ReadTimeout, cfg.WriteTimeout, zapLog)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zapLog.Fatalf("Ошибка сервера: %v", err)
	}

	zapLog.Infoln("Приложение завершило работу.")
}
