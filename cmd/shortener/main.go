package main

import (
	"errors"
	"log"
	"net/http"

	"mmskazak/shorturl/internal/app/config"
	"mmskazak/shorturl/internal/app/helpers"
	"mmskazak/shorturl/internal/pkg/shorturl"
)

func main() {
	appInfo, err := helpers.GetAppNameAndVersion()
	if err != nil {
		log.Printf("Ошибка при получении информации о приложении: %v", err)
	}

	log.Printf("Название приложения: %v", appInfo.Name)
	log.Printf("Версия: %v", appInfo.Version)

	cfg := config.InitConfig()

	app := shorturl.NewApp(cfg, cfg.ReadTimeout, cfg.WriteTimeout)

	if err := app.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
