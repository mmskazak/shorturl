package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mmskazak/shorturl/internal/app/helpers"
	"mmskazak/shorturl/internal/app/storage/mapstorage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type IStorage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

const (
	defaultShortURLLength = 8
	maxIteration          = 10
)

func CreateShortURL(w http.ResponseWriter, r *http.Request, data IStorage, baseHost string) {
	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read the body", http.StatusBadRequest)
		return
	}
	originalURL := string(body)

	// Генерируем уникальный идентификатор для сокращенной ссылки
	id := helpers.GenerateShortURL(defaultShortURLLength)

	err = data.SetShortURL(id, originalURL)
	if errors.Is(err, mapstorage.ErrKeyAlreadyExists) {
		var countIteration = 0
		for errors.Is(err, mapstorage.ErrKeyAlreadyExists) {
			id = helpers.GenerateShortURL(defaultShortURLLength)
			err = data.SetShortURL(id, originalURL)
			if countIteration == maxIteration {
				break
			}
			countIteration++
		}
	}

	shortedURL := baseHost + "/" + id

	if err != nil {
		errorGenerateShortURL := fmt.Errorf("ошибка формирования короткого url  %w", err)
		log.Printf("Ошибка SetShortURL: %v", errorGenerateShortURL)

		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("Сервису не удалось сформировать короткий URL"))
		if err != nil {
			log.Printf("Ошибка ответа ResponseWriter: %v ", err)
		}
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortedURL))
	if err != nil {
		log.Printf("Ошибка ResponseWriter %v", err)
	}
}

func HandleRedirect(w http.ResponseWriter, r *http.Request, data IStorage) {
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
