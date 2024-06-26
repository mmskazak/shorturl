package api

import (
	"context"
	"encoding/json"
	"errors"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"net/http"

	"go.uber.org/zap"
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
	zapLog *zap.SugaredLogger,
) {
	// Установка заголовков, чтобы указать, что мы принимаем и отправляем JSON
	w.Header().Set("Content-Type", "application/json")

	// Получаем userID из контекста
	payload, ok := r.Context().Value(ctxkeys.PayLoad).(jwtbuilder.PayloadJWT)
	userID := payload.UserID
	zapLog.Infof("FindUserURLs - UserID: %v", userID)
	if !ok {
		// Если userID не найден или неверного типа, возвращаем ошибку
		zapLog.Error("Не удалось получить id пользователя")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	// Получаем URL-адреса пользователя из базы данных
	urls, err := store.GetUserURLs(ctx, userID, baseHost)
	if errors.Is(err, storageErrors.ErrShortURLsForUserNotFound) {
		zapLog.Error("Не удалось получить URL-адреса пользователя из базы данных")
		http.Error(w, "", http.StatusNoContent)
		return
	}
	if err != nil {
		// Обработка ошибок, связанных с получением данных
		zapLog.Errorf("error getting user urls: %v", err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	// Преобразуем данные в JSON
	response, err := json.Marshal(urls)
	if err != nil {
		// Обработка ошибок, связанных с сериализацией JSON
		zapLog.Errorf("error marshalling user urls: %v", err)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	// Записываем JSON-ответ в ResponseWriter
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		zapLog.Errorf("error writing response: %v", err)
	}
}
