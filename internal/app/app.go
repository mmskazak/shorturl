package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/handlers/api"
	"mmskazak/shorturl/internal/handlers/web"
	"mmskazak/shorturl/internal/middleware"
	"mmskazak/shorturl/internal/storage"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

// Pinger определяет интерфейс для проверки состояния хранилища.
type Pinger interface {
	Ping(ctx context.Context) error
}

// URL представляет структуру URL с коротким и оригинальным URL.
type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// App представляет приложение с HTTP сервером и логгером.
type App struct {
	server *http.Server
	zapLog *zap.SugaredLogger
}

const ErrStartingServer = "error starting server"

// NewApp создает новый экземпляр приложения.
// Ctx - контекст для управления временем выполнения.
// Cfg - конфигурация приложения.
// Store - хранилище данных.
// ReadTimeout - таймаут чтения HTTP-запросов.
// WriteTimeout - таймаут записи HTTP-ответов.
// ZapLog - логгер.
func NewApp(
	ctx context.Context,
	cfg *config.Config,
	store storage.Storage,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	zapLog *zap.SugaredLogger,
) *App {
	router := chi.NewRouter()

	// Блок middleware
	router.Use(func(next http.Handler) http.Handler {
		return middleware.GetUserURLsForAuth(next, cfg)
	})
	router.Use(func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, cfg, zapLog)
	})
	router.Use(middleware.CheckUserID)
	router.Use(func(next http.Handler) http.Handler {
		return middleware.LoggingRequestMiddleware(next, zapLog)
	})
	router.Use(middleware.GzipMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		web.MainPage(w, r, zapLog)
	})

	baseHost := cfg.BaseHost // Получаем значение из конфига

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		zapLog.Infoln("Запрос получен handleRedirectHandler")
		web.HandleRedirect(ctx, w, r, store, zapLog)
	})

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		web.HandleCreateShortURL(ctx, w, r, store, baseHost, zapLog)
	})

	router.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		api.HandleCreateShortURL(ctx, w, r, store, baseHost, zapLog)
	})

	router.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		api.SaveShortenURLsBatch(ctx, w, r, store, cfg.BaseHost, zapLog)
	})

	pingPostgreSQL := func(w http.ResponseWriter, r *http.Request) {
		pinger, ok := store.(Pinger)
		if !ok {
			zapLog.Infoln("The storage does not support Ping")
			return
		}

		web.PingPostgreSQL(ctx, w, r, pinger, zapLog)
	}
	router.Get("/ping", pingPostgreSQL)

	router.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		api.FindUserURLs(ctx, w, r, store, cfg.BaseHost, zapLog)
	})

	router.Delete("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteUserURLs(ctx, w, r, store, zapLog)
	})

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
