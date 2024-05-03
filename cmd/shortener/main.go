package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"mmskazak/shorturl/internal/app/config"
	"mmskazak/shorturl/internal/app/helpers"
	"mmskazak/shorturl/internal/pkg/shorturl"
)

const readTimeout = 10 * time.Second
const writeTimeout = 10 * time.Second

// Точка входа программы.
func main() {
	appInfo := helpers.GetAppNameAndVersion()
	log.Println(appInfo)

	cfg := config.InitConfig()

	app := shorturl.NewApp(cfg, readTimeout, writeTimeout)

	if err := app.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
