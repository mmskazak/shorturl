package middleware

import (
	"net/http"

	"mmskazak/shorturl/internal/config"
)

func GetUserURLsForAuth(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := cfg.SecretKey
		_, err := getSignedPayloadJWT(r, authorizationCookieName, secretKey)
		if err != nil && r.URL.Path == "/api/user/urls" {
			w.WriteHeader(http.StatusUnauthorized)
		}
		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}
