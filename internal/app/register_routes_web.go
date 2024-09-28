package app

import (
	"context"
	"net/http"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/handlers/web"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func registrationWEBRoutes(
	ctx context.Context,
	router *chi.Mux,
	cfg *config.Config,
	zapLog *zap.SugaredLogger,
	store contracts.Storage,
	shortURLService contracts.IShortURLService,
) *chi.Mux {
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
		web.HandleCreateShortURL(ctx, w, r, store, baseHost, zapLog, shortURLService)
	})

	// Пинг PostgreSQL
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		pinger, ok := store.(contracts.Pinger)
		if !ok {
			zapLog.Infoln("The storage does not support Ping")
			return
		}

		web.PingPostgreSQL(ctx, w, r, pinger, zapLog)
	})

	return router
}
