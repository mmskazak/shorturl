package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"net/http"
)

type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func FindUserURLs(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	store storage.Storage,
	baseHost string,
) {
	// Установка заголовков, чтобы указать, что мы принимаем и отправляем JSON
	w.Header().Set("Content-Type", "application/json")

	// Получаем userID из контекста
	userID, ok := r.Context().Value(ctxkeys.KeyUserID).(string)
	if !ok {
		// Если userID не найден или неверного типа, возвращаем ошибку
		http.Error(w, "Не удалось получить id пользователя", http.StatusUnauthorized)
		return
	}

	// Получаем URL-адреса пользователя из базы данных
	urls, err := store.GetUserURLs(ctx, userID, baseHost)
	if errors.Is(err, storageErrors.ErrShortURLsForUserNotFound) {
		http.Error(w, "", http.StatusNoContent)
		return
	}
	if err != nil {
		// Обработка ошибок, связанных с получением данных
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	// Преобразуем данные в JSON
	response, err := json.Marshal(urls)
	if err != nil {
		// Обработка ошибок, связанных с сериализацией JSON
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	// Записываем JSON-ответ в ResponseWriter
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("error writing response: %v", err)
	}
}