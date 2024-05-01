package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"mmskazak/shorturl/internal/app/config"
	"mmskazak/shorturl/internal/app/handlers"
	"mmskazak/shorturl/internal/app/helpers"
	"mmskazak/shorturl/internal/app/middleware"

	"github.com/go-chi/chi/v5"
)

func main() {
	const timeoutDuration = 10 * time.Second

	app := helpers.GetAppNameAndVersion()
	log.Println(app)

	cfg := config.InitConfig()

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
