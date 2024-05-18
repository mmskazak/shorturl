package middleware

import (
	"compress/gzip"
	"mmskazak/shorturl/internal/logger"
	"net/http"
	"strings"
)

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
