package web

import (
	"context"
	"errors"
	"io"
	"log"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/shorturlservice"
	"mmskazak/shorturl/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type IGenIDForURL interface {
	Generate() (string, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

// Определяем тип для ключа контекста.
type contextKey string

// Постоянный ключ для идентификатора пользователя.
const keyUserID contextKey = "userID"

func HandleCreateShortURL(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	data storage.Storage,
	baseHost string) {
	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Не удалось прочитать тело запроса %v", err)
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}
	// Получаем userID из контекста
	userID, ok := r.Context().Value(keyUserID).(string)
	if !ok {
		// Если userID не найден или неверного типа, возвращаем ошибку
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	originalURL := string(body)
	generator := genidurl.NewGenIDService()
	shortURLService := shorturlservice.NewShortURLService()
	dto := shorturlservice.DTOShortURL{
		UserID:      userID,
		OriginalURL: originalURL,
		BaseHost:    baseHost,
	}

	shortURL, err := shortURLService.GenerateShortURL(ctx, dto, generator, data)
	if errors.Is(err, shorturlservice.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
		_, err := w.Write([]byte(shortURL))
		if err != nil {
			log.Printf("Ошибка записи ответа w.Write([]byte(shortURL)) при конфликте original url %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		return
	}

	if err != nil {
		log.Printf("Ошибка saveUniqueShortURL: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		log.Printf("Ошибка ResponseWriter: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func HandleRedirect(ctx context.Context, w http.ResponseWriter, r *http.Request, data storage.Storage) {
	// Получение значения id из URL-адреса
	id := chi.URLParam(r, "id")

	originalURL, err := data.GetShortURL(ctx, id)

	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func MainPage(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Сервис сокращения URL"))
	if err != nil {
		log.Printf("Ошибка при обращении к главной странице: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func PingPostgreSQL(ctx context.Context, w http.ResponseWriter, _ *http.Request, data Pinger) {
	err := data.Ping(ctx)
	if err != nil {
		log.Printf("Ошибка пинга базы данных data.Ping(ctx): %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
