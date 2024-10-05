package middleware

import (
	"net/http"

	"mmskazak/shorturl/internal/services/checkip"

	"go.uber.org/zap"
)

// IPRangeMiddleware мидлвар для проверки IP адреса по CIDR маске.
func IPRangeMiddleware(next http.Handler, cidr string, logger *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/internal/stats" {
			logger.Infoln("url path =", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}
		logger.Infoln("url path =", r.URL.Path)

		ip := realIP(r)
		logger.Infoln("ip =", ip)
		ok, err := checkip.CheckIPByCIDR(ip, cidr)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if !ok {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// realIP возвращает реальный IP адрес из заголовка X-Real-IP.
func realIP(r *http.Request) string {
	return r.Header.Get("X-Real-IP")
}
