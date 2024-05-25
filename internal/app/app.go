package app

import (
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/handlers/api"
	"mmskazak/shorturl/internal/handlers/web"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Storage второй раз объявляю интерфейс. Объявляем интерфейс где его используем.
type Storage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

type App struct {
	server *http.Server
}

const ErrStartingServer = "error starting server"

// NewApp создает новый экземпляр приложения.
func NewApp(cfg *config.Config, storage Storage, readTimeout time.Duration, writeTimeout time.Duration) *App {
	router := chi.NewRouter()

	// Добавление middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.GzipMiddleware)

	router.Get("/", web.MainPage)

	baseHost := cfg.BaseHost // Получаем значение из конфига

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	handleRedirectHandler := func(w http.ResponseWriter, r *http.Request) {
		web.HandleRedirect(w, r, storage)
	}
	router.Get("/{id}", handleRedirectHandler)

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	shortURLCreate := func(w http.ResponseWriter, r *http.Request) {
		web.HandleCreateShortURL(w, r, storage, baseHost)
	}
	router.Post("/", shortURLCreate)

	shortURLCreateAPI := func(w http.ResponseWriter, r *http.Request) {
		api.HandleCreateShortURL(w, r, storage, baseHost)
	}
	router.Post("/api/shorten", shortURLCreateAPI)

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
	logger.Log.Info("Server is running on %v", a.server.Addr)

	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Logf.Errorf("%v: %v", ErrStartingServer, err)
		return fmt.Errorf(ErrStartingServer+": %w", err)
	}
	return nil
}
