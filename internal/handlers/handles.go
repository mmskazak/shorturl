package handlers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"mmskazak/shorturl/config"
	"mmskazak/shorturl/internal/helpers"
	"mmskazak/shorturl/internal/repository"
	"net/http"
)

func CreateShortURL(w http.ResponseWriter, r *http.Request) {

	cfg := config.GetAppConfig()
	urlMap := repository.GetUrlMap()

	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read the body", http.StatusBadRequest)
		return
	}
	originalURL := string(body)

	// Генерируем уникальный идентификатор для сокращенной ссылки
	id := helpers.GenerateShortURL(8)
	shortedURL := cfg.BaseHost + "/" + id
	urlMap[id] = originalURL

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortedURL))
	if err != nil {
		return
	}
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	urlMap := repository.GetUrlMap()

	// Получение значения id из URL-адреса
	id := chi.URLParam(r, "id")

	if len(id) != 8 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL, ok := urlMap[id]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
