package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleRequests_GetUrl_BadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/wEr", nil)
	w := httptest.NewRecorder()

	handleRequests(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Ожидаетсся статус код %d, получен %d", http.StatusNotFound, res.StatusCode)
}

func TestHandleRequests_GetUrl_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/wErTdsBk", nil)
	w := httptest.NewRecorder()

	handleRequests(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode, "Ожидаетсся статус код %d, получен %d", http.StatusNotFound, res.StatusCode)
}

func TestHandleRequests_GetMainPage(t *testing.T) {
	//тестовый реквест
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder() //*ResponseRecorder

	handleRequests(w, req)

	//получаем результат *ResponseRecorder
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode, "Ожидаетсся статус код %d, получен %d", http.StatusOK, res.StatusCode)

	//сразу закрываем чтение
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	require.NoError(t, err, "Ошибка чтения(err) res.Body")

	assert.Equal(t, "Сервис сокращения URL", string(body))
}

func TestHandleRequests_PostMainPage(t *testing.T) {
	urlMap = make(map[string]string)

	originalURL := "https://ya.ru"
	requestBody := bytes.NewBuffer([]byte(originalURL))
	req := httptest.NewRequest(http.MethodPost, "/", requestBody)
	w := httptest.NewRecorder()

	handleRequests(w, req)

	res := w.Result()

	assert.Equal(t, http.StatusCreated, res.StatusCode, "Ожидаетсся статус код %d, получен %d", http.StatusCreated, res.StatusCode)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	require.NoError(t, err, "Ошибка чтения res.Body")

	shortenedURL := string(body)
	assert.NotEmpty(t, shortenedURL, "Ожидается сокращенный URL")

	// Additional checks if needed, e.g., checking if shortened URL is valid
}

func TestHandleRequests_GetUrl_Found(t *testing.T) {
	urlMap = make(map[string]string)
	urlMap["wErTdsBk"] = "https://ya.ru"

	req := httptest.NewRequest(http.MethodGet, "/wErTdsBk", nil)
	w := httptest.NewRecorder()

	handleRequests(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "Ожидаетсся статус код %d, получен %d", http.StatusTemporaryRedirect, res.StatusCode)
}
