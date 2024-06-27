package web

import (
	"context"
	"errors"
	"io"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"mmskazak/shorturl/internal/services/shorturlservice"
	"mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

type IGenIDForURL interface {
	Generate() (string, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

func HandleCreateShortURL(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	data storage.Storage,
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
	// Получаем userID из контекста
	payload, ok := r.Context().Value(ctxkeys.PayLoad).(jwtbuilder.PayloadJWT)
	userID := payload.UserID
	if !ok {
		zapLog.Infof("userID не найден или неверного типа, возвращаем http ошибку")
		// Если userID не найден или неверного типа, возвращаем ошибку
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	originalURL := string(body)
	generator := genidurl.NewGenIDService()
	shortURLService := shorturlservice.NewShortURLService()
	dto := shorturlservice.DTOShortURL{
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

func HandleRedirect(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	data storage.Storage,
	zapLog *zap.SugaredLogger,
) {
	// Получение значения id из URL-адреса
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

func MainPage(w http.ResponseWriter, _ *http.Request, zapLog *zap.SugaredLogger) {
	_, err := w.Write([]byte("Сервис сокращения URL"))
	if err != nil {
		zapLog.Errorf("Ошибка при обращении к главной странице: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func PingPostgreSQL(
	ctx context.Context,
	w http.ResponseWriter,
	_ *http.Request,
	data Pinger,
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
