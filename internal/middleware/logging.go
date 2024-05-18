package middleware

import (
	"fmt"
	"mmskazak/shorturl/internal/logger"
	"net/http"
	"time"
)

// Берём структуру для хранения сведений об ответе.
type responseData struct {
	status int
	size   int
}

// Добавляем реализацию http.ResponseWriter.
type loggingResponseWriter struct {
	ResponseWriter http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	responseData   *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return 0, fmt.Errorf("ошибка в middleware.loggingResponseWriter.Write %w", err)
	}
	r.responseData.size += size // захватываем размер
	return size, nil
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func (r *loggingResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

// LoggingMiddleware для логирования запросов.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		// Проход далее по цепочке middleware и обработчиков
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	})
}
