package middleware

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// LoggingMiddleware для логирования запросов.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логирование запроса
		log.Println("Incoming request:", r.Method, r.URL.Path)

		// Проход далее по цепочке middleware и обработчиков
		next.ServeHTTP(w, r)
	})
}
