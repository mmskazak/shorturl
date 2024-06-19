package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"net/http"
)

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
	userID, ok := r.Context().Value(ctxkeys.KeyUserID).(string)
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
