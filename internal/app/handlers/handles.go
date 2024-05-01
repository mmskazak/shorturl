package handlers

import (
	"io"
	"mmskazak/shorturl/internal/app/helpers"
	"mmskazak/shorturl/internal/app/storage/mapstorage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const defaultShortURLLength = 8

func CreateShortURL(w http.ResponseWriter, r *http.Request, baseHost string) {
	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read the body", http.StatusBadRequest)
		return
	}
	originalURL := string(body)

	// Генерируем уникальный идентификатор для сокращенной ссылки
	id := helpers.GenerateShortURL(defaultShortURLLength)
	shortedURL := baseHost + "/" + id
	data := mapstorage.GetMapStorageInstance()

	err = data.SetShortURL(id, originalURL)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortedURL))
	if err != nil {
		return
	}
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	data := mapstorage.GetMapStorageInstance()

	// Получение значения id из URL-адреса
	id := chi.URLParam(r, "id")

	if len(id) != defaultShortURLLength {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL, err := data.GetShortURL(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func MainPage(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Сервис сокращения URL"))
	if err != nil {
		return
	}
}
