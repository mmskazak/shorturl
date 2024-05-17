package app

import (
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/handlers"
	"mmskazak/shorturl/internal/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// IStorage второй раз объявляю интерфейс. Объявляем интерфейс где его используем.
type IStorage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

type App struct {
	server *http.Server
}

const ErrStartingServer = "error starting server"

// NewApp создает новый экземпляр приложения.
func NewApp(cfg *config.Config, storage IStorage, readTimeout time.Duration, writeTimeout time.Duration) *App {
	router := chi.NewRouter()

	// Добавление middleware
	router.Use(middleware.LoggingMiddleware)

	router.Get("/", handlers.MainPage)

	baseHost := cfg.BaseHost // Получаем значение из конфига

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	handleRedirectHandler := func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRedirect(w, r, storage)
	}
	router.Get("/{id}", handleRedirectHandler)

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	createShortURLHandler := func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateShortURL(w, r, storage, baseHost)
	}
	router.Post("/", createShortURLHandler)

	return &App{
		server: &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

// Start запускает сервер приложения.
func (a *App) Start() error {
	log.Printf("Server is running on %v", a.server.Addr)

	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("%v: %v", ErrStartingServer, err)
		return fmt.Errorf(ErrStartingServer+": %w", err)
	}
	return nil
}
