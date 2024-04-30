package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"mmskazak/shorturl/config"
	"mmskazak/shorturl/internal/handlers"
	"mmskazak/shorturl/internal/middleware"
	mapstorage "mmskazak/shorturl/internal/storage/map-storage"
	"net/http"
	"os"
)

var cfg *config.Config

func init() {
	// Создание нового экземпляра конфигурации
	cfg = config.InitConfig()
	_ = mapstorage.GetMapStorageInstance()
}

func main() {

	// делаем разбор командной строки
	flag.Parse()

	//конфигурационные параметры в приоритете из переменных среды
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

	fmt.Println("Server is running on " + cfg.Address)
	err := http.ListenAndServe(cfg.Address, router)
	if err != nil {
		return
	}
}
