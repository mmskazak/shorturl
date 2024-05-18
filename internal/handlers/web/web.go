package web

import (
	"io"
	"log"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/shorturlservice"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func HandleCreateShortURL(w http.ResponseWriter, r *http.Request, storage IStorage, baseHost string) {
	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Не удалось прочитать тело запроса %v", err)
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}
	originalURL := string(body)
	generator := genidurl.NewGenIDService()
	shortURLService := shorturlservice.NewShortURLService()
	dto := shorturlservice.DTOShortURL{
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
