package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/shorturlservice"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"net/http"
)

type JSONRequest struct {
	URL string `json:"url"`
}

type JSONResponse struct {
	ShortURL string `json:"result"`
}

// Определяем тип для ключа контекста.
type contextKey string

const (
	appJSON = "application/json"
	// Постоянный ключ для идентификатора пользователя.
	keyUserID contextKey = "userID"
)

func HandleCreateShortURL(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	store storage.Storage,
	baseHost string,
) {
	// Установка заголовков, чтобы указать, что мы принимаем и отправляем JSON.
	w.Header().Set("Content-Type", appJSON)
	w.Header().Set("Accept", appJSON)

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

	jsonReq := JSONRequest{}

	err = json.Unmarshal(body, &jsonReq)
	if err != nil {
		log.Printf("Ошибка json.Unmarshal: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	generator := genidurl.NewGenIDService()
	shortURLService := shorturlservice.NewShortURLService()
	dto := shorturlservice.DTOShortURL{
		UserID:      userID,
		OriginalURL: jsonReq.URL,
		BaseHost:    baseHost,
	}

	shortURL, err := shortURLService.GenerateShortURL(ctx, dto, generator, store)
	if errors.Is(err, shorturlservice.ErrConflict) {
		shortURLAsJSON, err := buildJSONResponse(shortURL)
		if err != nil {
			log.Printf("Ошибка buildJSONResponse: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte(shortURLAsJSON))
		if err != nil {
			log.Printf("Ошибка write, err := w.Write([]byte(shortURLAsJson)): %v", err)
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
	jsonResp := JSONResponse{
		ShortURL: shortURL,
	}
	shortURLAsJSON, err := json.Marshal(jsonResp)
	if err != nil {
		log.Printf("Ошибка json.Marshal: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(shortURLAsJSON)
	if err != nil {
		log.Printf("Ошибка ResponseWriter: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func SaveShortenURLsBatch(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	store storage.Storage,
	baseHost string,
) {
	// Парсинг JSON из тела запроса
	var requestData []storage.Incoming
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID, ok := r.Context().Value(keyUserID).(string)
	if !ok {
		// Если userID не найден или неверного типа, возвращаем ошибку
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	generator := genidurl.NewGenIDService()
	// Сохранение пакета коротких URL
	outputs, err := store.SaveBatch(ctx, requestData, baseHost, userID, generator)
	if err != nil {
		if errors.Is(err, storageErrors.ErrUniqueViolation) {
			http.Error(w, "", http.StatusConflict)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Преобразование результата в JSON
	responseData, err := json.Marshal(outputs)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа
	w.Header().Set("Content-Type", appJSON)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func buildJSONResponse(shortURL string) (string, error) {
	jsonResp := JSONResponse{
		ShortURL: shortURL,
	}
	shortURLAsJSON, err := json.Marshal(jsonResp)
	if err != nil {
		return "", fmt.Errorf("ошибка json marshal %w", err)
	}

	return string(shortURLAsJSON), nil
}
