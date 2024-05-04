package main

import (
	"mmskazak/shorturl/internal/app/config"
	"mmskazak/shorturl/internal/app/handlers"
	"mmskazak/shorturl/internal/app/storage/mapstorage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestMainPage_Get_Greeting(t *testing.T) {
	// Создаем фейковый роутер chi с нашим обработчиком
	r := chi.NewRouter()

	r.Get("/", handlers.MainPage)

	// Создаем фейковый HTTP запрос
	req, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Вызываем обработчик с фейковым запросом
	r.ServeHTTP(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := "Сервис сокращения URL"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateShortURL_Post_Create(t *testing.T) {
	r := chi.NewRouter()

	ms := mapstorage.NewMapStorage()
	cfg := config.InitConfig()
	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	createShortURLHandler := func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateShortURL(w, r, ms, cfg.GetBaseHost())
	}
	r.Post("/", createShortURLHandler)

	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
}

func TestHandleRedirect_Get(t *testing.T) {
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
			ms := mapstorage.NewMapStorage()
			err := ms.SetShortURL("vAlIdIds", "https://ya.ru")
			if err != nil {
				t.Fatal(err)
			}

			handleRedirectHandler := func(w http.ResponseWriter, r *http.Request) {
				handlers.HandleRedirect(w, r, ms)
			}
			r.Get("/{id}", handleRedirectHandler)

			req, err := http.NewRequest(http.MethodGet, tc.path, http.NoBody)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
		})
	}
}
