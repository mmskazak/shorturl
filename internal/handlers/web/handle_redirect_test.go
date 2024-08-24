package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"mmskazak/shorturl/internal/contracts/mocks"

	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"github.com/go-chi/chi/v5"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestHandleRedirectWithRouterChi(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/IDShtURL", http.NoBody)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockIGetShortURL(ctrl)
	data.EXPECT().GetShortURL(gomock.Any(), "IDShtURL").Return("http://yandex.ru", nil)

	// Настройка роутера chi для обработки запроса
	router := chi.NewRouter()
	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		HandleRedirect(ctx, w, r, data, zapSugar)
	})

	// Выполнение запроса через роутер
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestHandleRedirect(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/IDShtURL", http.NoBody)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockIGetShortURL(ctrl)
	data.EXPECT().GetShortURL(gomock.Any(), gomock.Any()).Return("http://yandex.ru", nil)

	HandleRedirect(context.Background(), w, req, data, zapSugar)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestHandleRedirect_ErrDeletedShortURL(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/IDShtURL", http.NoBody)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockIGetShortURL(ctrl)
	data.EXPECT().GetShortURL(gomock.Any(), gomock.Any()).Return("", storageErrors.ErrDeletedShortURL)

	HandleRedirect(context.Background(), w, req, data, zapSugar)

	assert.Equal(t, http.StatusGone, w.Code)
}

func TestHandleRedirect_ErrGetShortURL(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/IDShtURL", http.NoBody)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockIGetShortURL(ctrl)
	data.EXPECT().GetShortURL(gomock.Any(), gomock.Any()).Return("", errors.New("test error"))

	HandleRedirect(context.Background(), w, req, data, zapSugar)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
