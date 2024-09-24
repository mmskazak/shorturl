package middleware

import (
	"fmt"
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

		ipNet, err := parseCIDR(cidr)
		if err != nil {
			logger.Infof("error parsing cidr %q: %v", cidr, err)
			http.Error(w, "", http.StatusForbidden)
			return
		}

		ip := realIP(r)
		if ip != "" {
			clientIP := net.ParseIP(ip)
			if clientIP != nil && ipNet.Contains(clientIP) {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "", http.StatusForbidden)
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
