package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/dtos"
	"mmskazak/shorturl/internal/models"

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
	shortURLService contracts.IShortURLService,
) {
	// Установка заголовков, чтобы указать, что мы принимаем и отправляем JSON.
	w.Header().Set("Content-Type", appJSON)
	w.Header().Set("Accept", appJSON)

	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
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
	zapLog.Infof("User ID: %d", userID)
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
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte(shortURLAsJSON))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		return
	}
	zapLog.Infoln("Short URL: ", shortURL)

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
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	zapLog.Infoln("Short URL JSON: ", string(shortURLAsJSON))
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(shortURLAsJSON)
	if err != nil {
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
