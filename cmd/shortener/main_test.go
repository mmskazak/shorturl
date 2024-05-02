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
	req, err := http.NewRequest("GET", "/", nil)
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
		handlers.CreateShortURL(w, r, ms, cfg.BaseHost)
	}
	r.Post("/", createShortURLHandler)

	req, err := http.NewRequest("POST", "/", strings.NewReader("https://ya.ru"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
}

func TestHandleRedirect_Get_Found(t *testing.T) {
	r := chi.NewRouter()

	ms := mapstorage.NewMapStorage()
	err := ms.SetShortURL("vAlIdIds", "https://ya.ru")
	if err != nil {
		t.Fatal(err)
	}

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	handleRedirectHandler := func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRedirect(w, r, ms)
	}
	r.Get("/{id}", handleRedirectHandler)

	req, err := http.NewRequest("GET", "/vAlIdIds", strings.NewReader("https://ya.ru")) //nolint:usestdlibvars
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
}

func TestHandleRedirect_Get_NotFound(t *testing.T) {
	r := chi.NewRouter()
	ms := mapstorage.NewMapStorage()

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	handleRedirectHandler := func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRedirect(w, r, ms)
	}
	r.Get("/{id}", handleRedirectHandler)

	req, err := http.NewRequest("GET", "/x0x0x0x0", strings.NewReader("https://ya.ru")) //nolint:usestdlibvars
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound)
}

func TestHandleRedirect_Get_BadRequest(t *testing.T) {
	r := chi.NewRouter()
	ms := mapstorage.NewMapStorage()

	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
	handleRedirectHandler := func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRedirect(w, r, ms)
	}
	r.Get("/{id}", handleRedirectHandler)

	req, err := http.NewRequest("GET", "/x0x0", strings.NewReader("https://ya.ru")) //nolint:usestdlibvars
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
