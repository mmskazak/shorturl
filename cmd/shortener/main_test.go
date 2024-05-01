package main

//import (
//	"bytes"
//	"mmskazak/shorturl/internal/app/handlers"
//	"mmskazak/shorturl/internal/app/storage/mapstorage"
//
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/go-chi/chi/v5"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)
//
//func TestMainPage_Get_Greeting(t *testing.T) {
//	// тестовый реквест
//	req := httptest.NewRequest(http.MethodGet, "/", nil)
//	w := httptest.NewRecorder() // *ResponseRecorder
//
//	handlers.MainPage(w, req)
//
//	// получаем результат *ResponseRecorder
//	res := w.Result()
//	assert.Equal(t, http.StatusOK, res.StatusCode, "Ожидаетсся статус код %d, получен %d", http.StatusOK, res.StatusCode)
//
//	// сразу закрываем чтение
//	defer res.Body.Close()
//	body, err := io.ReadAll(res.Body)
//
//	require.NoError(t, err, "Ошибка чтения(err) res.Body")
//
//	assert.Equal(t, "Сервис сокращения URL", string(body))
//}
//
//func TestCreateShortURL_Post_Create(t *testing.T) {
//	originalURL := "https://ya.ru"
//	requestBody := bytes.NewBufferString(originalURL)
//	req := httptest.NewRequest(http.MethodPost, "/", requestBody)
//	w := httptest.NewRecorder()
//
//	// Создаем тестовый маршрутизатор Chi
//	r := chi.NewRouter()
//
//	// Создаем замыкание, которое передает значение конфига в обработчик CreateShortURL
//	createShortURLHandler := func(w http.ResponseWriter, r *http.Request) {
//		baseHost := cfg.BaseHost // Получаем значение из конфига
//		handlers.CreateShortURL(w, r, baseHost)
//	}
//	// Определяем маршрут для обработки функции handleRedirect
//	r.Post("/{id}", createShortURLHandler)
//
//	//handlers.CreateShortURL(w, req,)
//
//	res := w.Result()
//
//	assert.Equal(
//		t,
//		http.StatusCreated,
//		res.StatusCode,
//		"Ожидается статус код %d, получен %d",
//		http.StatusCreated,
//		res.StatusCode,
//	)
//
//	defer res.Body.Close()
//	body, err := io.ReadAll(res.Body)
//
//	require.NoError(t, err, "Ошибка чтения res.Body")
//
//	shortenedURL := string(body)
//	assert.NotEmpty(t, shortenedURL, "Ожидается сокращенный URL")
//}
//
//func TestHandleRedirect_Get_Found(t *testing.T) {
//	ms := mapstorage.NewMapStorage()
//	mapstorage.SetMapStorageInstance(ms)
//
//	err := ms.SetShortURL("validID8", "https://ya.ru")
//	if err != nil {
//		return
//	}
//
//	// Создаем тестовый маршрутизатор Chi
//	r := chi.NewRouter()
//
//	// Определяем маршрут для обработки функции handleRedirect
//	r.Get("/{id}", handlers.HandleRedirect)
//
//	// Создаем тестовый запрос с URL-параметром "id"
//	req := httptest.NewRequest(http.MethodGet, "/validID8", nil)
//
//	// Создаем тестовый ответ (используется только для записи результата)
//	w := httptest.NewRecorder()
//
//	// Обработаем тестовый запрос с помощью тестового маршрутизатора
//	r.ServeHTTP(w, req)
//
//	// Проверим статус код ответа
//	res := w.Result()
//	defer res.Body.Close()
//	assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "Ожидается статус код 307")
//
//	// Проверяем фактический URL перенаправления
//	assert.Equal(t, "https://ya.ru", res.Header.Get("Location"), "Ожидается перенаправление на https://ya.ru")
//}
//
//func TestHandleRedirect_Get_NotFound(t *testing.T) {
//	ms := mapstorage.NewMapStorage()
//	mapstorage.SetMapStorageInstance(ms)
//
//	err := ms.SetShortURL("validIDx", "https://ya.ru")
//	if err != nil {
//		return
//	}
//
//	// Создаем тестовый маршрутизатор Chi
//	r := chi.NewRouter()
//
//	// Определяем маршрут для обработки функции handleRedirect
//	r.Get("/{id}", handlers.HandleRedirect)
//
//	// Создаем тестовый запрос с несуществующим ID
//	req := httptest.NewRequest(http.MethodGet, "/validID8", nil)
//
//	// Создаем тестовый ответ (используется только для записи результата)
//	w := httptest.NewRecorder()
//
//	// Обработаем тестовый запрос с помощью тестового маршрутизатора
//	r.ServeHTTP(w, req)
//
//	// Проверим статус код ответа - Not Found (404)
//	res := w.Result()
//	defer res.Body.Close()
//	assert.Equal(t, http.StatusNotFound, res.StatusCode, "Ожидается статус код 404")
//}
//
//func TestHandleRedirect_Get_BadRequest(t *testing.T) {
//	ms := mapstorage.GetMapStorageInstance()
//	err := ms.SetShortURL("validID", "https://ya.ru")
//	if err != nil {
//		return
//	}
//
//	// Создаем тестовый маршрутизатор Chi
//	r := chi.NewRouter()
//
//	// Определяем маршрут для обработки функции handleRedirect
//	r.Get("/{id}", handlers.HandleRedirect)
//
//	// Создаем тестовый запрос с неправильным форматом ID (слишком короткий)
//	req := httptest.NewRequest(http.MethodGet, "/NotValidID", nil)
//
//	// Создаем тестовый ответ (используется только для записи результата)
//	w := httptest.NewRecorder()
//
//	// Обработайте тестовый запрос с помощью тестового маршрутизатора
//	r.ServeHTTP(w, req)
//
//	// Проверим статус код ответа - Bad Request (400)
//	res := w.Result()
//	defer res.Body.Close()
//	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Ожидается статус код 400")
//}
