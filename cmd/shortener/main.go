package main

import (
	"flag"
	"log"
	"mmskazak/shorturl/config"
	"mmskazak/shorturl/internal/handlers"
	"mmskazak/shorturl/internal/helpers"
	"mmskazak/shorturl/internal/middleware"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

var cfg *config.Config //nolint:gochecknoglobals

func main() {
	app := helpers.GetAppNameAndVersion()
	log.Println(app)

	cfg = config.InitConfig()
	// делаем разбор командной строки
	flag.Parse()

	// конфигурационные параметры в приоритете из переменных среды
	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		cfg.Address = envServAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseHost = envBaseURL
	}

	router := chi.NewRouter()

	// Добавление middleware
	router.Use(middleware.LoggingMiddleware)

	router.Get("/", handlers.MainPage)
	router.Get("/{id}", handlers.HandleRedirect)
	router.Post("/", handlers.CreateShortURL)

	log.Println("Server is running on " + cfg.Address)
	var err = http.ListenAndServe(cfg.Address, router)
	if err != nil {
		return
	}
}
