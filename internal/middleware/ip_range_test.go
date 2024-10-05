package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestIPRangeMiddleware(t *testing.T) {
	// Моковый обработчик, который будет вызван, если IP допустим
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := zap.NewNop().Sugar() // Мокаем логгер

	t.Run("valid IP in CIDR range", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", http.NoBody)
		req.Header.Set("X-Real-IP", "192.168.0.1") // Устанавливаем IP, который будет проверяться

		rr := httptest.NewRecorder()

		middleware := IPRangeMiddleware(nextHandler, "192.168.0.0/24", logger)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid IP out of CIDR range", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", http.NoBody)
		req.Header.Set("X-Real-IP", "10.0.0.1") // IP, не входящий в CIDR

		rr := httptest.NewRecorder()

		middleware := IPRangeMiddleware(nextHandler, "192.168.0.0/24", logger)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("non-stats endpoint should pass through", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/other-endpoint", http.NoBody)
		req.Header.Set("X-Real-IP", "192.168.0.1")

		rr := httptest.NewRecorder()

		middleware := IPRangeMiddleware(nextHandler, "192.168.0.0/24", logger)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("check IP by CIDR error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", http.NoBody)
		req.Header.Set("X-Real-IP", "192.168.0.1")

		rr := httptest.NewRecorder()

		middleware := IPRangeMiddleware(nextHandler, "errorCIDR", logger)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
