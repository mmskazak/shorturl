package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/dtos"
	"mmskazak/shorturl/internal/models"
	"net/http"

	"go.uber.org/zap"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"mmskazak/shorturl/internal/services/shorturlservice"
)

const (
	appJSON = "application/json"
)

// HandleCreateShortURL обрабатывает HTTP-запрос для создания короткого URL.
// Он принимает JSON-запрос с оригинальным URL, генерирует короткий URL и сохраняет его в хранилище.
// Если короткий URL уже существует, возвращает конфликт.
func HandleCreateShortURL(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	store contracts.ISetShortURL,
	baseHost string,
	zapLog *zap.SugaredLogger,
) {
	// Установка заголовков, чтобы указать, что мы принимаем и отправляем JSON.
	w.Header().Set("Content-Type", appJSON)
	w.Header().Set("Accept", appJSON)

	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		zapLog.Errorf("Не удалось прочитать тело запроса %v", err)
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		zapLog.Error("Тело запроса пустое, ошибка.")
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	payload, ok := r.Context().Value(ctxkeys.PayLoad).(jwtbuilder.PayloadJWT)
	userID := payload.UserID
	if !ok {
		// Если userID не найден или неверного типа, возвращаем ошибку
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	jsonReq := models.JSONRequest{}
	err = json.Unmarshal(body, &jsonReq)
	if err != nil {
		zapLog.Errorf("Ошибка json.Unmarshal: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	generator := genidurl.NewGenIDService()
	shortURLService := shorturlservice.NewShortURLService()
	dto := dtos.DTOShortURL{
		UserID:      userID,
		OriginalURL: jsonReq.URL,
		BaseHost:    baseHost,
		Deleted:     false,
	}

	shortURL, err := shortURLService.GenerateShortURL(ctx, dto, generator, store)
	if errors.Is(err, shorturlservice.ErrConflict) {
		shortURLAsJSON, err := buildJSONResponse(shortURL)
		if err != nil {
			zapLog.Errorf("Ошибка buildJSONResponse: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte(shortURLAsJSON))
		if err != nil {
			zapLog.Errorf("Ошибка write, err := w.Write([]byte(shortURLAsJson)): %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		return
	}

	if err != nil {
		zapLog.Errorf("Ошибка saveUniqueShortURL: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	jsonResp := models.JSONResponse{
		ShortURL: shortURL,
	}
	shortURLAsJSON, err := json.Marshal(jsonResp)
	if err != nil {
		zapLog.Errorf("Ошибка json.Marshal: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(shortURLAsJSON)
	if err != nil {
		zapLog.Errorf("Ошибка ResponseWriter: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// buildJSONResponse создает JSON-ответ с коротким URL.
func buildJSONResponse(shortURL string) (string, error) {
	jsonResp := models.JSONResponse{
		ShortURL: shortURL,
	}
	shortURLAsJSON, err := json.Marshal(jsonResp)
	if err != nil {
		return "", fmt.Errorf("ошибка json marshal %w", err)
	}

	return string(shortURLAsJSON), nil
}
