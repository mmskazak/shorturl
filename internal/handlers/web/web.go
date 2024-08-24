package web

import (
	"context"
	"errors"
	"io"
	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/dtos"
	"net/http"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"mmskazak/shorturl/internal/services/shorturlservice"
	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

//go:generate mockgen -source=web.go -destination=mocks/mock_web.go -package=mocks

// HandleCreateShortURL обрабатывает запрос на создание короткого URL.
// Он извлекает оригинальный URL из тела запроса, генерирует короткий URL и сохраняет его в хранилище.
// Возвращает HTTP-ответ с созданным коротким URL или ошибку в случае неудачи.
func HandleCreateShortURL(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	data contracts.ISetShortURL,
	baseHost string,
	zapLog *zap.SugaredLogger,
) {
	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		zapLog.Errorf("Не удалось прочитать тело запроса %v", err)
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		zapLog.Error("Тело запроса пустое")
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста.
	payload, ok := r.Context().Value(ctxkeys.PayLoad).(jwtbuilder.PayloadJWT)
	userID := payload.UserID
	if !ok {
		zapLog.Infof("userID не найден или неверного типа, возвращаем http ошибку")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	originalURL := string(body)
	generator := genidurl.NewGenIDService()
	shortURLService := shorturlservice.NewShortURLService()
	dto := dtos.DTOShortURL{
		UserID:      userID,
		OriginalURL: originalURL,
		BaseHost:    baseHost,
		Deleted:     false,
	}

	shortURL, err := shortURLService.GenerateShortURL(ctx, dto, generator, data)
	if errors.Is(err, shorturlservice.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
		_, err := w.Write([]byte(shortURL))
		if err != nil {
			zapLog.Errorf("Ошибка записи ответа w.Write([]byte(shortURL)) при конфликте original url %v", err)
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

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		zapLog.Errorf("Ошибка ResponseWriter: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// HandleRedirect обрабатывает запрос на перенаправление по короткому URL.
// Получает короткий URL из параметра запроса, извлекает оригинальный URL из хранилища и выполняет перенаправление.
// Возвращает HTTP-ответ с кодом перенаправления или ошибку в случае неудачи.
func HandleRedirect(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	data contracts.IGetShortURL,
	zapLog *zap.SugaredLogger,
) {
	// Получение значения id из URL-адреса.
	id := chi.URLParam(r, "id")

	originalURL, err := data.GetShortURL(ctx, id)
	if errors.Is(err, storageErrors.ErrDeletedShortURL) {
		zapLog.Errorf("error is deleted shorturl: %v", err)
		http.Error(w, "", http.StatusGone)
		return
	}

	if err != nil {
		zapLog.Errorf("error is getting shorturl: %v", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

// MainPage обрабатывает запросы к главной странице.
// Возвращает сообщение о службе сокращения URL или ошибку в случае неудачи.
func MainPage(w http.ResponseWriter, _ *http.Request, zapLog *zap.SugaredLogger) {
	_, err := w.Write([]byte("Сервис сокращения URL"))
	if err != nil {
		zapLog.Errorf("Ошибка при обращении к главной странице: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// PingPostgreSQL проверяет состояние подключения к базе данных PostgreSQL.
// Возвращает HTTP-ответ с кодом состояния 200 OK, если база данных доступна, или ошибку в случае неудачи.
func PingPostgreSQL(
	ctx context.Context,
	w http.ResponseWriter,
	_ *http.Request,
	data contracts.Pinger,
	zapLog *zap.SugaredLogger,
) {
	err := data.Ping(ctx)
	if err != nil {
		zapLog.Errorf("Ошибка пинга базы данных data.Ping(ctx): %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
