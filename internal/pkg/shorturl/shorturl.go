package shorturl

import (
	"github.com/go-chi/chi/v5"
	"log"
	"mmskazak/shorturl/internal/app/config"
	"mmskazak/shorturl/internal/app/handlers"
	"mmskazak/shorturl/internal/app/middleware"
	"net/http"
	"time"
)

// почему линтен это не отлеживает?
const timeoutDuration = 10 * time.Second

type App struct {
	Config *config.Config
	Router *chi.Mux
	Server *http.Server
}

// NewApp создает новый экземпляр приложения
func NewApp(cfg *config.Config, readTimeout time.Duration, writeTimeout time.Duration) *App {
	r := chi.NewRouter()

	// Добавление middleware
	r.Use(middleware.LoggingMiddleware)

	r.Get("/", handlers.MainPage)
	r.Get("/{id}", handlers.HandleRedirect)
	r.Post("/", handlers.CreateShortURL)

	return &App{
		Config: cfg,
		Router: r,
		Server: &http.Server{
			Addr:         cfg.Address,
			Handler:      r,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

// Start запускает сервер приложения
func (a *App) Start() error {
	log.Println("Server is running on " + a.Config.Address)
	return a.Server.ListenAndServe()
}
