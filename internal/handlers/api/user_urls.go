package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"mmskazak/shorturl/internal/contracts"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"go.uber.org/zap"
)

//go:generate mockgen -source=user_urls.go -destination=mocks/mock_user_urls.go -package=mocks

// FindUserURLs обрабатывает запрос на получение всех URL, созданных пользователем.
// Он извлекает userID из контекста, получает соответствующие URL из хранилища и возвращает их клиенту в формате JSON.
// Если возникают ошибки, возвращаются соответствующие HTTP-статус коды и сообщения об ошибке.
func FindUserURLs(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	store contracts.IGetUserURLs,
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
