package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckUserID(t *testing.T) {
	// Создаем HTTP тестовый сервер с применением middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем статус ответа
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	middleware := CheckUserID(nextHandler)

	// Тестовый случай 1: JWT отсутствует или недействителен
	req, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	middleware.ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
