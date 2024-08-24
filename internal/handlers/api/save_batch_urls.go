package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/models"
	"net/http"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"go.uber.org/zap"
)

//go:generate mockgen -source=save_batch_urls.go -destination=mocks/mock_save_batch_urls.go -package=mocks

// SaveShortenURLsBatch обрабатывает пакетный запрос на создание сокращённых URL.
// Он парсит JSON-тело запроса, извлекает userID из контекста,
// сохраняет пакет сокращённых URL и возвращает результат клиенту в виде JSON.
// Если возникает ошибка при парсинге, сохранении или маршалинге,
// возвращается соответствующий HTTP-статус код и сообщение об ошибке.
func SaveShortenURLsBatch(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	store contracts.ISaveBatch,
	baseHost string,
	zapLog *zap.SugaredLogger,
) {
	// Парсинг JSON из тела запроса
	var requestData []models.Incoming
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		zapLog.Errorf("error decoding request: %v", err)
		http.Error(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	payload, ok := r.Context().Value(ctxkeys.PayLoad).(jwtbuilder.PayloadJWT)
	userID := payload.UserID
	if !ok {
		// Если userID не найден или неверного типа, возвращаем ошибку
		zapLog.Error("error getting user id from context")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	generator := genidurl.NewGenIDService()
	// Сохранение пакета коротких URL
	outputs, err := store.SaveBatch(ctx, requestData, baseHost, userID, generator)
	if err != nil {
		if errors.Is(err, storageErrors.ErrUniqueViolation) {
			zapLog.Errorw("error saving shorten URLs batch", "error", err)
			http.Error(w, "", http.StatusConflict)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Преобразование результата в JSON
	responseData, err := json.Marshal(outputs)
	if err != nil {
		zapLog.Errorw("error marshalling shorten URLs batch", "error", err)
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
