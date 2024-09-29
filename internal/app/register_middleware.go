package app

import (
	"net/http"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/middleware"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// registrationMiddleware регистрация мидлваров.
func registrationMiddleware(router *chi.Mux, cfg *config.Config, zapLog *zap.SugaredLogger) *chi.Mux {
	// Блок проверки IP адреса по CIDR маске
	router.Use(func(next http.Handler) http.Handler {
		zapLog.Infoln("middleware IPRangeMiddleware проверки IP адреса по CIDR маске")
		return middleware.IPRangeMiddleware(next, cfg.TrustedSubnet, zapLog)
	})
	// Блок middleware
	router.Use(func(next http.Handler) http.Handler {
		zapLog.Infoln("middleware GetUserURLsForAuth для запросов к пути /api/user/urls")
		return middleware.GetUserURLsForAuth(next, cfg)
	})
	router.Use(func(next http.Handler) http.Handler {
		zapLog.Infoln("middleware AuthMiddleware для аутентификации запросов на основе JWT")
		return middleware.AuthMiddleware(next, cfg, zapLog)
	})
	router.Use(middleware.CheckUserID)
	router.Use(func(next http.Handler) http.Handler {
		zapLog.Infoln("middleware LoggingRequestMiddleware логирует информацию о каждом HTTP-запросе и ответе")
		return middleware.LoggingRequestMiddleware(next, zapLog)
	})
	router.Use(middleware.GzipMiddleware)

	return router
}
