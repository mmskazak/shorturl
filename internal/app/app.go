package app

import (
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/handlers/api"
	"mmskazak/shorturl/internal/handlers/web"
	"mmskazak/shorturl/internal/middleware"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

type Storage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

type App struct {
	server *http.Server
}

const ErrStartingServer = "error starting server"

// NewApp создает новый экземпляр приложения.
func NewApp(cfg *config.Config,
	storage Storage,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	zapLog *zap.SugaredLogger) *App {
	router := chi.NewRouter()

	// Add the custom logging middleware to the router
	LoggingMiddlewareRich := func(next http.Handler) http.Handler {
		return middleware.LoggingMiddleware(next, zapLog)
	}

	// Добавление middleware
	router.Use(LoggingMiddlewareRich)
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
	log.Printf("Server is running on %v\n", a.server.Addr)

	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("%v: %v", ErrStartingServer, err)
		return fmt.Errorf(ErrStartingServer+": %w", err)
	}
	return nil
}
