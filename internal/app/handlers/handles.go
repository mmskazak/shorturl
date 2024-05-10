package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mmskazak/shorturl/internal/app/helpers"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type IStorage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

const (
	defaultShortURLLength  = 8
	maxIteration           = 10
	InternalServerErrorMsg = "Внутренняя ошибка сервера"
)

func CreateShortURL(w http.ResponseWriter, r *http.Request, storage IStorage, baseHost string) {
	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Не удалось прочитать тело запроса %v", err)
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}
	originalURL := string(body)

	id, err := saveUniqueShortURL(storage, originalURL)

	if err != nil {
		log.Printf("Ошибка saveUniqueShortURL: %v", err)
		http.Error(w, "Сервису не удалось сформировать короткий URL", http.StatusInternalServerError)
		return
	}

	base, err := url.Parse(baseHost)
	if err != nil {
		log.Printf("Ошибка при разборе базового URL: %v", err)
		http.Error(w, InternalServerErrorMsg, http.StatusInternalServerError)
		return
	}

	idPath, err := url.Parse(id)
	if err != nil {
		log.Printf("Ошибка при разборе пути ID: %v", err)
		http.Error(w, InternalServerErrorMsg, http.StatusInternalServerError)
		return
	}

	shortURL := base.ResolveReference(idPath)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL.String()))
	if err != nil {
		log.Printf("Ошибка ResponseWriter: %v", err)
		http.Error(w, InternalServerErrorMsg, http.StatusInternalServerError)
		return
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
		http.Error(w, InternalServerErrorMsg, http.StatusInternalServerError)
		log.Printf("Ошибка при обращении к главной странице: %v", err)
	}
}

func saveUniqueShortURL(storage IStorage, originalURL string) (string, error) {
	for range maxIteration {
		id, err := helpers.GenerateShortURL(defaultShortURLLength)
		if err != nil {
			return "", fmt.Errorf("функция GenerateShortURL вернула ошибку %w", err)
		}

		err = storage.SetShortURL(id, originalURL)
		if err == nil {
			return id, nil
		} else {
			return "", fmt.Errorf("метод SetShortURL вернула ошибку %w", err)
		}
	}
	return "", errors.New("не удалось сгенерировать уникальный идентификатор")
}
