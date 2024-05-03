package handlers

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"mmskazak/shorturl/internal/app/storage/mapstorage"
)

func TestCreateShortURL(t *testing.T) {
	// Initialize a new MapStorage for testing
	data := mapstorage.NewMapStorage()

	// Define a test handler function that wraps CreateShortURL
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateShortURL(w, r, data, "https://ya.ru")
	})

	// Create a new request with a POST method and a body containing the original URL
	reqBody := bytes.NewBufferString("https://ya.ru")
	req := httptest.NewRequest(http.MethodPost, "/", reqBody)

	// Create a new response recorder to capture the response from the handler
	w := httptest.NewRecorder()

	// Call the handler function with the test request and response recorder
	handler.ServeHTTP(w, req)

	// Check that the response status code is 201 Created
	require.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

func TestHandleRedirect(t *testing.T) {
	// Initialize a new MapStorage for testing
	data := *mapstorage.NewMapStorage()
	err := data.SetShortURL("x5x5x5x5", "https://ya.ru")
	require.NoError(t, err)

	// Create a new chi router
	r := chi.NewRouter()

	// Define the route and bind the handler
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		HandleRedirect(w, r, &data)
	})

	// Создаем фейковый HTTP запрос
	req, err := http.NewRequest(http.MethodGet, "/x5x5x5x5", nil)
	require.NoError(t, err)
	// Create a new response recorder to capture the response from the handler
	w := httptest.NewRecorder()

	// Call the handler function with the test request and response recorder
	r.ServeHTTP(w, req)

	// Check that the response status code is 201 Created
	require.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

func TestMainPage(t *testing.T) {
	// Create a new chi router
	r := chi.NewRouter()

	// Define the route and bind the handler
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r)
	})

	// Создаем фейковый HTTP запрос
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	// Create a new response recorder to capture the response from the handler
	w := httptest.NewRecorder()

	// Call the handler function with the test request and response recorder
	r.ServeHTTP(w, req)

	// Check that the response status code is 201 Created
	require.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, w.Body.String(), "Сервис сокращения URL")

}
