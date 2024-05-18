package middleware

import (
	"compress/gzip"
	"fmt"
	"mmskazak/shorturl/internal/logger"
	"net/http"
	"strings"
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

// Добавляем реализацию http.ResponseWriter.
type gzipResponseWriter struct {
	ResponseWriter http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	gzWriter       *gzip.Writer
}

// Write - переопределенный метод Write для сжатия данных.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	// Записываем данные в сжатый поток с использованием gzip.Writer
	countByte, err := w.gzWriter.Write(b)
	if err != nil {
		logger.Log.Errorln(err)
		return 0, fmt.Errorf("переопределенный метод Write для сжатия данных, ошибка: %w", err)
	}
	return countByte, nil
}

// Header - реализация метода Header интерфейса http.ResponseWriter.
func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// WriteHeader - реализация метода WriteHeader интерфейса http.ResponseWriter.
func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
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

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Log.Errorln(err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
			defer func(gz *gzip.Reader) {
				err := gz.Close()
				if err != nil {
					logger.Log.Errorln(err)
				}
			}(gz)

			gzWriter := gzip.NewWriter(w)
			defer func(gzWriter *gzip.Writer) {
				err := gzWriter.Close()
				if err != nil {
					logger.Log.Errorln(err)
				}
			}(gzWriter)

			gzResponseWriter := &gzipResponseWriter{
				gzWriter:       gzWriter,
				ResponseWriter: w,
			}

			w.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(gzResponseWriter, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
