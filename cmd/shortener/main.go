package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"

	"mmskazak/shorturl/internal/app"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/storage/factory"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
	storage, err := factory.NewStorage(ctx, cfg, zapLog)
	if err != nil {
		zapLog.Fatalf("Ошибка инициализации хранилища: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			zapLog.Warn("Error closing storage: %v\n", err)
		}
	}()

	newApp := app.NewApp(ctx, cfg, storage, cfg.ReadTimeout, cfg.WriteTimeout, zapLog)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zapLog.Fatalf("Ошибка сервера: %v", err)
	}

	zapLog.Infoln("Приложение завершило работу.")
}
