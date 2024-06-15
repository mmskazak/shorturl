package api

import (
	"context"
	"encoding/json"
	"log"
	"mmskazak/shorturl/internal/storage"
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
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		// Если userID не найден или неверного типа, возвращаем ошибку
		http.Error(w, "Не удалось получить id пользователя", http.StatusInternalServerError)
		return
	}

	// Получаем URL-адреса пользователя из базы данных
	urls, err := store.GetUserURLs(ctx, userID, baseHost)
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
