package shorturl

import (
	"log"
	"mmskazak/shorturl/internal/app/config"
	"mmskazak/shorturl/internal/app/handlers"
	"mmskazak/shorturl/internal/app/middleware"
	"mmskazak/shorturl/internal/app/storage/mapstorage"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	server *http.Server
}

// NewApp создает новый экземпляр приложения.
func NewApp(cfg *config.Config, readTimeout time.Duration, writeTimeout time.Duration) *App {
	router := chi.NewRouter()

	// Добавление middleware
	router.Use(middleware.LoggingMiddleware)

	router.Get("/", handlers.MainPage)

	ms := mapstorage.NewMapStorage()
	baseHost := cfg.BaseHost // Получаем значение из конфига

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	handleRedirectHandler := func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRedirect(w, r, ms)
	}
	router.Get("/{id}", handleRedirectHandler)

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	createShortURLHandler := func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateShortURL(w, r, ms, baseHost)
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
	log.Println("Server is running on " + a.server.Addr)
	return a.server.ListenAndServe()
}
