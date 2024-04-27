package main

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage_Get_Greeting(t *testing.T) {
	//тестовый реквест
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder() //*ResponseRecorder

	mainPage(w, req)

	//получаем результат *ResponseRecorder
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode, "Ожидаетсся статус код %d, получен %d", http.StatusOK, res.StatusCode)

	//сразу закрываем чтение
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	require.NoError(t, err, "Ошибка чтения(err) res.Body")

	assert.Equal(t, "Сервис сокращения URL", string(body))
}

func TestCreateShortURL_Post_Create(t *testing.T) {
	urlMap = make(map[string]string)

	originalURL := "https://ya.ru"
	requestBody := bytes.NewBuffer([]byte(originalURL))
	req := httptest.NewRequest(http.MethodPost, "/", requestBody)
	w := httptest.NewRecorder()

	createShortURL(w, req)

	res := w.Result()

	assert.Equal(t, http.StatusCreated, res.StatusCode, "Ожидаетсся статус код %d, получен %d", http.StatusCreated, res.StatusCode)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	require.NoError(t, err, "Ошибка чтения res.Body")

	shortenedURL := string(body)
	assert.NotEmpty(t, shortenedURL, "Ожидается сокращенный URL")
}

func TestHandleRedirect_Get_Found(t *testing.T) {
	urlMap := make(map[string]string)
	urlMap["validID8"] = "https://ya.ru"

	// Создаем тестовый запрос с URL-параметром "id"
	req := httptest.NewRequest(http.MethodGet, "/validID8", nil)

	// Создаем тестовый ответ (используется только для записи результата)
	w := httptest.NewRecorder()

	// Обрабатываем тестовый запрос с помощью тестового маршрутизатора
	handleRedirect(w, req)

	// Проверьте статус код ответа
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "Ожидается статус код 307")
}

func TestHandleRedirect_Get_NotFound(t *testing.T) {
	urlMap := make(map[string]string)
	urlMap["validIDx"] = "https://ya.ru"

	// Создайте тестовый маршрутизатор Chi
	r := chi.NewRouter()

	// Определите маршрут для обработки функции handleRedirect
	r.Get("/{id}", handleRedirect)

	// Создайте тестовый запрос с несуществующим ID
	req := httptest.NewRequest(http.MethodGet, "/invalid1", nil)

	// Создайте тестовый ответ (используется только для записи результата)
	w := httptest.NewRecorder()

	// Обработайте тестовый запрос с помощью тестового маршрутизатора
	r.ServeHTTP(w, req)

	// Проверьте статус код ответа - Not Found (404)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode, "Ожидается статус код 404")
}

func TestHandleRedirect_Get_BadRequest(t *testing.T) {
	urlMap := make(map[string]string)
	urlMap["validID"] = "https://ya.ru"

	// Создайте тестовый маршрутизатор Chi
	r := chi.NewRouter()

	// Определите маршрут для обработки функции handleRedirect
	r.Get("/{id}", handleRedirect)

	// Создайте тестовый запрос с неправильным форматом ID (слишком короткий)
	req := httptest.NewRequest(http.MethodGet, "/short", nil)

	// Создайте тестовый ответ (используется только для записи результата)
	w := httptest.NewRecorder()

	// Обработайте тестовый запрос с помощью тестового маршрутизатора
	r.ServeHTTP(w, req)

	// Проверьте статус код ответа - Bad Request (400)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Ожидается статус код 400")
}
