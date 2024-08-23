package api

import (
	"context"
	"errors"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/handlers/api/mocks"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"mmskazak/shorturl/internal/storage"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestFindUserURLs_Success(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Создание контекста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req = req.WithContext(ctx)

	// Базовый хост
	baseHost := "http://localhost"

	expectedOutput := []storage.URL{
		{
			ShortURL:    "http://localhost:8080/eythhwGV",
			OriginalURL: "https://google.ru",
		},
	}

	ctxBg := context.Background()

	// Создание мока для ISetShortURL
	data := mocks.NewMockIGetUserURLs(ctrl)
	data.EXPECT().GetUserURLs(ctxBg, "11111", baseHost).Return(expectedOutput, nil)

	FindUserURLs(ctxBg, w, req, data, baseHost, zapSugar)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"short_url\":\"http://localhost:8080/eythhwGV\","+
		"\"original_url\":\"https://google.ru\"}]",
		w.Body.String())
}

func TestFindUserURLs_StatusUnauthorized(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctxBg := context.Background()

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	req = req.WithContext(ctxBg)

	// Базовый хост
	baseHost := "http://localhost"

	// Создание мока для ISetShortURL
	data := mocks.NewMockIGetUserURLs(ctrl)

	FindUserURLs(ctxBg, w, req, data, baseHost, zapSugar)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFindUserURLs_StatusOK(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем контекст реквеста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создаем обычный контекст
	ctxBg := context.Background()

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	req = req.WithContext(ctx)

	// Базовый хост
	baseHost := "http://localhost"

	// Создание мока для ISetShortURL
	data := mocks.NewMockIGetUserURLs(ctrl)
	expectedURL := []storage.URL{
		{
			ShortURL:    "https://localhost:8080/nQm6WEim",
			OriginalURL: "https://yandex.ru",
		},
	}
	data.EXPECT().GetUserURLs(ctxBg, "11111", baseHost).Return(expectedURL, nil)

	FindUserURLs(ctxBg, w, req, data, baseHost, zapSugar)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"short_url\":\"https://localhost:8080/nQm6WEim\","+
		"\"original_url\":\"https://yandex.ru\"}]",
		w.Body.String())
}

func TestFindUserURLs_StatusNoContent(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем контекст реквеста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создаем обычный контекст
	ctxBg := context.Background()

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	req = req.WithContext(ctx)

	// Базовый хост
	baseHost := "http://localhost"

	// Создание мока для ISetShortURL
	data := mocks.NewMockIGetUserURLs(ctrl)

	data.EXPECT().GetUserURLs(ctxBg, "11111", baseHost).
		Return(nil, storageErrors.ErrShortURLsForUserNotFound)

	FindUserURLs(ctxBg, w, req, data, baseHost, zapSugar)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestFindUserURLs_StatusInternalServerError(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем контекст реквеста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создаем обычный контекст
	ctxBg := context.Background()

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	req = req.WithContext(ctx)

	// Базовый хост
	baseHost := "http://localhost"

	// Создание мока для ISetShortURL
	data := mocks.NewMockIGetUserURLs(ctrl)

	data.EXPECT().GetUserURLs(ctxBg, "11111", baseHost).
		Return(nil, errors.New("test error"))

	FindUserURLs(ctxBg, w, req, data, baseHost, zapSugar)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
