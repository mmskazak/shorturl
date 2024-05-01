package main

import (
	"errors"
	"flag"
	"log"
	"mmskazak/shorturl/config"
	"mmskazak/shorturl/internal/handlers"
	"mmskazak/shorturl/internal/helpers"
	"mmskazak/shorturl/internal/middleware"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	const timeoutDuration = 10 * time.Second

	app := helpers.GetAppNameAndVersion()
	log.Println(app)

	cfg := config.InitConfig()
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

	// Создаем сервер
	srv := &http.Server{
		Addr:         cfg.Address,     // cfg.Address - адрес сервера из конфигурации
		Handler:      router,          // router - HTTP маршрутизатор
		ReadTimeout:  timeoutDuration, // Время ожидания на чтение запроса
		WriteTimeout: timeoutDuration, // Время ожидания на запись ответа
	}

	// Запускаем сервер с явным указанием параметров таймаута
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
