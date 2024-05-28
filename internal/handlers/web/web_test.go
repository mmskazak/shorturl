package web

import (
	"bytes"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage/infile"
	"mmskazak/shorturl/internal/storage/inmemory"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage(t *testing.T) {
	// Create a new chi router
	r := chi.NewRouter()

	// Define the route and bind the handler
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r)
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

func TestCreateShortURL(t *testing.T) {
	// Initialize a new MapStorage for testing
	cfg := config.Config{
		FileStoragePath: "/tmp/file.json",
	}
	ms, _ := infile.NewInFile(&cfg)

	// Define a test handler function that wraps CreateShortURL
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCreateShortURL(w, r, ms, "http://ya.ru")
	})

	// Create a new request with a POST method and a body containing the original URL
	reqBody := bytes.NewBufferString("http://ya.ru")
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
	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{
			name:         "NotFound",
			path:         "/x0x0x0x0",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "BadRequest",
			path:         "/x0x0",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Redirect",
			path:         "/vAlIdIds",
			expectedCode: http.StatusTemporaryRedirect,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := chi.NewRouter()
			ms, err := inmemory.NewInMemory()
			require.NoError(t, err)
			err = ms.SetShortURL("vAlIdIds", "http://ya.ru")
			require.NoError(t, err)

			handleRedirectHandler := func(w http.ResponseWriter, r *http.Request) {
				HandleRedirect(w, r, ms)
			}
			r.Get("/{id}", handleRedirectHandler)

			req, err := http.NewRequest(http.MethodGet, tc.path, http.NoBody)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
		})
	}
}
