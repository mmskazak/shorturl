package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	writer io.Writer
	http.ResponseWriter
}

// Write переопределяет метод Write для GzipResponseWriter,
// чтобы записывать данные в gzip.Writer вместо обычного ResponseWriter.
func (w *GzipResponseWriter) Write(b []byte) (int, error) {
	write, err := w.writer.Write(b)
	if err != nil {
		return 0, fmt.Errorf("error writing to gzip writer: %w", err)
	}
	return write, nil
}

// GzipMiddleware обрабатывает запросы и ответы с поддержкой сжатия gzip.
// Если запрос содержит сжатое тело (gzip), оно будет декомпрессировано.
// Если ответ должен быть сжат (gzip), он будет сжат перед отправкой клиенту.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Декомпрессия тела запроса, если оно сжато (gzip)
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Invalid gzip body", http.StatusBadRequest)
				return
			}
			defer func(gz *gzip.Reader) {
				err := gz.Close()
				if err != nil {
					log.Printf("error gzReader close %v", err)
				}
			}(gzipReader)
			r.Body = gzipReader
		}

		// Проверка, нужно ли сжимать ответ
		contentType := w.Header().Get("Content-Type")
		isCompressingContent := strings.HasPrefix(contentType, "application/json") ||
			strings.HasPrefix(contentType, "text/html")

		// Сжатие ответа (gzip), если клиент поддерживает gzip и контент подходящего типа
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && isCompressingContent {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")

			gzipWriter := gzip.NewWriter(w)
			defer func(gzipWriter *gzip.Writer) {
				err := gzipWriter.Close()
				if err != nil {
					log.Printf("error close gzipWriter %v", err)
				}
			}(gzipWriter)

			// Оборачиваем ResponseWriter в GzipResponseWriter для сжатия ответа
			gzipResponseWriter := &GzipResponseWriter{writer: gzipWriter, ResponseWriter: w}
			next.ServeHTTP(gzipResponseWriter, r)
			return
		}
		// Если gzip не требуется, просто передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}
