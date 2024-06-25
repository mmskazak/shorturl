package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage(t *testing.T) {
	// Create a new chi router
	r := chi.NewRouter()

	// Define the route and bind the handler
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger, err := zap.NewDevelopment()
		require.NoError(t, err)
		zapLog := logger.Sugar()
		MainPage(w, r, zapLog)
	})

	// Создаем фейковый HTTP запрос
	req, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	require.NoError(t, err)

	// Create a new response recorder to capture the response from the handler
	w := httptest.NewRecorder()

	// Call the handler function with the test request and response recorder
	r.ServeHTTP(w, req)

	// Check that the response status code is 201 Created
	require.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, w.Body.String(), "Сервис сокращения URL")
}
