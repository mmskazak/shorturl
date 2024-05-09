package main

import (
	"errors"
	"log"
	"mmskazak/shorturl/internal/app"
	"mmskazak/shorturl/internal/app/storage/mapstorage"
	"net/http"

	"mmskazak/shorturl/internal/app/config"
	"mmskazak/shorturl/internal/app/helpers"
)

func main() {

	appInfo, err := helpers.GetAppNameAndVersion()
	if err != nil {
		log.Printf("Ошибка при получении информации о приложении: %v", err)
	} else {
		log.Printf("Название приложения: %v", appInfo.Name)
		log.Printf("Версия: %v", appInfo.Version)
	}

	cfg, errCfg := config.InitConfig()
	if err != nil {
		log.Printf("ошибка инициализации конфигурации в main %v", errCfg)
	}

	ms := mapstorage.NewMapStorage()

	newApp := app.NewApp(cfg, ms, cfg.ReadTimeout, cfg.WriteTimeout)

	if err := newApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
