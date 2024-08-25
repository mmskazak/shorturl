package middleware

import (
	"mmskazak/shorturl/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserURLsForAuth(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "secret",
	}

	// Создаем фейковый обработчик для проверки результата
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := GetUserURLsForAuth(handler, cfg)

	tests := []struct {
		name           string
		cookieValue    string
		requestPath    string
		expectedStatus int
	}{
		{
			name: "Valid JWT on /api/user/urls",
			cookieValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJ1c2VyX2lkIjoiODUyOGQ2MDYtZjEyZi00MTdkLTk3YjAtNDljOWE3NjE2M2FhIn0." +
				"IAoOF6UEjDuK5BIMoSisZEoaGyB8Yb4Z_hdn75_3nu0",
			requestPath:    "/api/user/urls",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JWT on /api/user/urls",
			cookieValue:    "invalid_token",
			requestPath:    "/api/user/urls",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "No JWT on /api/user/urls",
			cookieValue:    "",
			requestPath:    "/api/user/urls",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "No JWT on other path",
			cookieValue:    "",
			requestPath:    "/api/other/path",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.requestPath, http.NoBody)
			require.NoError(t, err)

			// Устанавливаем JWT в куки, если это необходимо
			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{
					Name:  authorizationCookieName,
					Value: tt.cookieValue,
				})
			}

			rr := httptest.NewRecorder()
			middleware.ServeHTTP(rr, req)
			// Проверяем статус ответа
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
