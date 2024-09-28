package app

import (
	"context"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/handlers/api"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func registrationAPIRoutes(
	ctx context.Context,
	router *chi.Mux,
	cfg *config.Config,
	zapLog *zap.SugaredLogger,
	store contracts.Storage,
	shortURLService contracts.IShortURLService,
) *chi.Mux {
	baseHost := cfg.BaseHost // Получаем значение из конфига

	router.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		api.HandleCreateShortURL(ctx, w, r, store, baseHost, zapLog, shortURLService)
	})

	router.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		api.SaveShortenURLsBatch(ctx, w, r, store, cfg.BaseHost, zapLog)
	})

	router.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		api.FindUserURLs(ctx, w, r, store, cfg.BaseHost, zapLog)
	})

	router.Delete("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteUserURLs(ctx, w, r, store, zapLog)
	})

	router.Get("/api/internal/stats", func(w http.ResponseWriter, r *http.Request) {
		api.InternalStats(ctx, w, r, store, zapLog)
	})

	return router
}
