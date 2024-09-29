package middleware

import (
	"net/http"

	"mmskazak/shorturl/internal/services/checkip"

	"go.uber.org/zap"
)

// IPRangeMiddleware мидлвар для проверки IP адреса по CIDR маске.
func IPRangeMiddleware(next http.Handler, cidr string, logger *zap.SugaredLogger) http.Handler {
	logger = logger.Named("IPRangeMiddleware")
	logger.Debug("IPRangeMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/internal/stats" {
			logger.Infoln("url path =", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}

		ip := realIP(r)
		ok, err := checkip.CheckIPByCIDR(ip, cidr)
		if err != nil {
			logger.Errorw("Error checking ip range", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if !ok {
			logger.Errorw("Error checking ip range", "ip", ip, "cidr", cidr)
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
