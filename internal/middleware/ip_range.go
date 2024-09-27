package middleware

import (
	"fmt"
	"mmskazak/shorturl/internal/services/checkip"
	"net"
	"net/http"

	"go.uber.org/zap"
)

// IPRangeMiddleware мидлвар для проверки IP адреса по CIDR маске.
func IPRangeMiddleware(next http.Handler, cidr string, logger *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/internal/stats" {
			next.ServeHTTP(w, r)
			return
		}

		ip := realIP(r)
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

// parseCIDR парсит CIDR строку и возвращает объект net.IPNet.
func parseCIDR(cidr string) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("parse cidr %s error: %w", cidr, err)
	}
	return ipNet, nil
}

// realIP возвращает реальный IP адрес из заголовка X-Real-IP.
func realIP(r *http.Request) string {
	return r.Header.Get("X-Real-IP")
}
