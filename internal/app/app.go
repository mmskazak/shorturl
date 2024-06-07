package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/handlers/api"
	"mmskazak/shorturl/internal/handlers/web"
	"mmskazak/shorturl/internal/middleware"
	storageInterface "mmskazak/shorturl/internal/storage"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type App struct {
	server *http.Server
}

const ErrStartingServer = "error starting server"

// NewApp создает новый экземпляр приложения.
func NewApp(
	ctx context.Context,
	cfg *config.Config,
	data storageInterface.Storage,
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
		zapLog.Infoln("Запрос получен handleRedirectHandler")
		web.HandleRedirect(ctx, w, r, data)
	}
	router.Get("/{id}", handleRedirectHandler)

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	handleCreateShortURL := func(w http.ResponseWriter, r *http.Request) {
		web.HandleCreateShortURL(ctx, w, r, data, baseHost)
	}
	router.Post("/", handleCreateShortURL)

	shortURLCreateAPI := func(w http.ResponseWriter, r *http.Request) {
		api.HandleCreateShortURL(ctx, w, r, data, baseHost)
	}
	router.Post("/api/shorten", shortURLCreateAPI)

	handleSaveShortenURLsBatch := func(w http.ResponseWriter, r *http.Request) {
		api.SaveShortenURLsBatch(ctx, w, r, data, cfg.BaseHost)
	}
	router.Post("/api/shorten/batch", handleSaveShortenURLsBatch)

	pingPostgreSQL := func(w http.ResponseWriter, r *http.Request) {
		pinger, ok := data.(Pinger)
		if !ok {
			http.Error(w, ErrStartingServer, http.StatusInternalServerError)
			return
		}

		web.PingPostgreSQL(ctx, w, r, pinger)
	}
	router.Get("/ping", pingPostgreSQL)

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
