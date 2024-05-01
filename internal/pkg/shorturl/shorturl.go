package shorturl

import (
	"log"
	"mmskazak/shorturl/internal/app/config"
	"mmskazak/shorturl/internal/app/handlers"
	"mmskazak/shorturl/internal/app/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	config *config.Config
	router *chi.Mux
	server *http.Server
}

// NewApp создает новый экземпляр приложения.
func NewApp(cfg *config.Config, readTimeout time.Duration, writeTimeout time.Duration) *App {
	router := chi.NewRouter()

	// Добавление middleware
	router.Use(middleware.LoggingMiddleware)

	router.Get("/", handlers.MainPage)
	router.Get("/{id}", handlers.HandleRedirect)

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	createShortURLHandler := func(w http.ResponseWriter, r *http.Request) {
		baseHost := cfg.BaseHost // Получаем значение из конфига
		handlers.CreateShortURL(w, r, baseHost)
	}
	router.Post("/", createShortURLHandler)

	return &App{
		config: cfg,
		router: router,
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
	log.Println("Server is running on " + a.config.Address)
	return a.server.ListenAndServe()
}
