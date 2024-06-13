package app

import (
	"context"
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/handlers/api"
	"mmskazak/shorturl/internal/handlers/web"
	"mmskazak/shorturl/internal/middleware"
	"mmskazak/shorturl/internal/storage"
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
	zapLog *zap.SugaredLogger
}

const ErrStartingServer = "error starting server"

// NewApp создает новый экземпляр приложения.
func NewApp(
	ctx context.Context,
	cfg *config.Config,
	store storage.Storage,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	zapLog *zap.SugaredLogger,
) *App {
	router := chi.NewRouter()

	// Add the custom logging middleware to the router
	LoggingMiddlewareRich := func(next http.Handler) http.Handler {
		return middleware.LoggingMiddleware(next, zapLog)
	}

	// Добавление middleware
	router.Use(middleware.AuthMiddleware)
	router.Use(LoggingMiddlewareRich)
	router.Use(middleware.GzipMiddleware)

	router.Get("/", web.MainPage)

	baseHost := cfg.BaseHost // Получаем значение из конфига

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		zapLog.Infoln("Запрос получен handleRedirectHandler")
		web.HandleRedirect(ctx, w, r, store)
	})

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		web.HandleCreateShortURL(ctx, w, r, store, baseHost)
	})

	router.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		api.HandleCreateShortURL(ctx, w, r, store, baseHost)
	})

	router.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		api.SaveShortenURLsBatch(ctx, w, r, store, cfg.BaseHost)
	})

	pingPostgreSQL := func(w http.ResponseWriter, r *http.Request) {
		pinger, ok := store.(Pinger)
		if !ok {
			zapLog.Infoln("The storage does not support Ping")
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
		zapLog: zapLog,
	}
}

// Start запускает сервер приложения.
func (a *App) Start() error {
	a.zapLog.Infof("Server is running on %v\n", a.server.Addr)

	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.zapLog.Infof("%v: %v", ErrStartingServer, err)
		return fmt.Errorf(ErrStartingServer+": %w", err)
	}
	return nil
}
