package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"mmskazak/shorturl/internal/app/storage/mapstorage"
)

func TestCreateShortURL(t *testing.T) {
	t.Skip()
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
	t.Skip()
	// Initialize a new MapStorage for testing
	data := *mapstorage.NewMapStorage()
	err := data.SetShortURL("x5x5x5x5", "https://ya.ru")
	require.NoError(t, err)

	// Define a test handler function that wraps CreateShortURL
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleRedirect(w, r, &data)
	})

	// Создаем фейковый HTTP запрос
	req, err := http.NewRequest(http.MethodGet, "/x5x5x5x5", nil)
	require.NoError(t, err)
	// Create a new response recorder to capture the response from the handler
	w := httptest.NewRecorder()

	// Call the handler function with the test request and response recorder
	handler.ServeHTTP(w, req)

	// Check that the response status code is 201 Created
	require.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

func TestMainPage(t *testing.T) {
	t.Skip()
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success main page",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MainPage(tt.args.w, tt.args.req)
		})
	}
}
