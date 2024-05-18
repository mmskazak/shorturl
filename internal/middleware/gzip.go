package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"mmskazak/shorturl/internal/logger"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	Writer io.Writer
	http.ResponseWriter
}

func (w *GzipResponseWriter) Write(b []byte) (int, error) {
	write, err := w.Writer.Write(b)
	if err != nil {
		return 0, fmt.Errorf("error writing to gzip writer: %w", err)
	}
	return write, nil
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle gzip request body
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Invalid gzip body", http.StatusBadRequest)
				return
			}
			defer func(gz *gzip.Reader) {
				err := gz.Close()
				if err != nil {
					logger.Log.Errorln(err)
				}
			}(gzipReader)
			r.Body = gzipReader
		}

		contentType := w.Header().Get("Content-Type")
		isCompressingContent := strings.HasPrefix(contentType, "application/json") ||
			strings.HasPrefix(contentType, "text/html")

		// Handle gzip response
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && isCompressingContent {
			gzipWriter := gzip.NewWriter(w)
			defer func(gzipWriter *gzip.Writer) {
				err := gzipWriter.Close()
				if err != nil {
					logger.Log.Errorln(err)
				}
			}(gzipWriter)

			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")
			gzipResponseWriter := &GzipResponseWriter{Writer: gzipWriter, ResponseWriter: w}
			next.ServeHTTP(gzipResponseWriter, r)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}