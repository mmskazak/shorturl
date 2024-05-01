package middleware

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// Conveyor
//
//	func main() {
//		http.Handle("/", Conveyor(http.HandlerFunc(rootHandle), middleware1, middleware2, middleware3))
//		// ...
//	}
//
// почему линтер не увидел не используемую структуру?
func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// LoggingMiddleware Middleware для логирования запросов.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логирование запроса
		log.Println("Incoming request:", r.Method, r.URL.Path)

		// Проход далее по цепочке middleware и обработчиков
		next.ServeHTTP(w, r)
	})
}
