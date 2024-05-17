package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"mmskazak/shorturl/internal/services"
	"net/http"
)

type IStorage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

type IGenIDForURL interface {
	Generate(int) (string, error)
}

const (
	defaultShortURLLength  = 8
	maxIteration           = 10
	InternalServerErrorMsg = "Внутренняя ошибка сервера"
)

var ErrOriginalURLIsEmpty = errors.New("originalURL is empty")
var ErrMethodSetShortURLNotCanSave = errors.New("метод SetShortURL не смог сохранить URL, ошибка")
var ErrFuncGenerate = errors.New("функция GenerateShortURL вернула ошибку")

func CreateShortURL(w http.ResponseWriter, r *http.Request, storage IStorage, baseHost string) {
	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Не удалось прочитать тело запроса %v", err)
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}
	originalURL := string(body)
	generator := services.NewGenIDService()
	shortURLService := services.NewShortURLService()
	dto := services.DTOShortURL{
		OriginalURL:  originalURL,
		BaseHost:     baseHost,
		MaxIteration: maxIteration,
		LengthID:     defaultShortURLLength,
	}

	shortURL, err := shortURLService.GenerateShortURL(dto, generator, storage)
	if err != nil {
		log.Printf("Ошибка saveUniqueShortURL: %v", err)
		http.Error(w, "Сервису не удалось сформировать короткий URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
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

func saveUniqueShortURL(storage IStorage, generator IGenIDForURL, originalURL string) (string, error) {
	if originalURL == "" {
		return "", ErrOriginalURLIsEmpty
	}

	var err error
	for range maxIteration {
		id, err := generator.Generate(defaultShortURLLength)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrFuncGenerate, err)
		}

		err = storage.SetShortURL(id, originalURL)
		if err == nil {
			return id, nil
		}
	}
	return "", fmt.Errorf("%w: %w", ErrMethodSetShortURLNotCanSave, err)
}
