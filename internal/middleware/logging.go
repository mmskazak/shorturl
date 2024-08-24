package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	noStatus = 0 // Константа для обозначения отсутствия статуса
)

// responseData хранит информацию о статусе ответа и размере данных.
type responseData struct {
	status int // HTTP статус-код ответа
	size   int // Размер ответа в байтах
}

// loggingResponseWriter реализует http.ResponseWriter и добавляет функциональность для логирования.
type loggingResponseWriter struct {
	http.ResponseWriter               // Встраиваем оригинальный http.ResponseWriter
	responseData        *responseData // Ссылка на структуру для хранения данных о ответе
}

// Write переопределяет метод Write для записи ответа и сбора размера ответа.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return 0, fmt.Errorf("ошибка в middleware.loggingResponseWriter.Write: %w", err)
	}
	r.responseData.size += size // Увеличиваем размер ответа
	if r.responseData.status == noStatus {
		r.responseData.status = http.StatusOK // Устанавливаем статус по умолчанию, если он ещё не был установлен
	}
	return size, nil
}

// WriteHeader переопределяет метод WriteHeader для записи кодов статуса ответа и их хранения.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // Сохраняем статус ответа
}

// LoggingRequestMiddleware логирует информацию о каждом HTTP-запросе и ответе.
func LoggingRequestMiddleware(next http.Handler, zapLog *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Засекаем время начала обработки запроса

		responseData := &responseData{
			status: noStatus, // Инициализируем статус как нет установленного значения
			size:   0,        // Инициализируем размер как ноль
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		// Передаем запрос следующему обработчику
		next.ServeHTTP(&lw, r)

		duration := time.Since(start) // Вычисляем длительность обработки запроса

		// Логируем информацию о запросе и ответе
		zapLog.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status, // Логируем статус ответа
			"duration", duration, // Логируем длительность обработки запроса
			"size", responseData.size, // Логируем размер ответа
		)
	})
}
