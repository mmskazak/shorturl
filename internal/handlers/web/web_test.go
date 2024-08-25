package web

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"mmskazak/shorturl/internal/contracts/mocks"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"mmskazak/shorturl/internal/services/shorturlservice"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleCreateShortURL_ErrConflict(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()
	// Создание HTTP-запроса с телом запроса
	body := strings.NewReader(`{"url": "https://google.ru"}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockISetShortURL(ctrl)

	// Базовый хост
	baseHost := "http://localhost"

	shortURLService := mocks.NewMockIShortURLService(ctrl)
	shortURLService.EXPECT().GenerateShortURL(context.Background(), gomock.Any(), gomock.Any(), data).
		Return("", shorturlservice.ErrConflict)

	// Вызов функции
	HandleCreateShortURL(context.Background(), w, req, data, baseHost, zapSugar, shortURLService)

	// Проверка кода ответа
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestHandleCreateShortURL_ErrGenerateShortURL(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "11111"})

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()
	// Создание HTTP-запроса с телом запроса
	body := strings.NewReader(`{"url": "https://google.ru"}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockISetShortURL(ctrl)

	// Базовый хост
	baseHost := "http://localhost"

	shortURLService := mocks.NewMockIShortURLService(ctrl)
	shortURLService.EXPECT().GenerateShortURL(context.Background(), gomock.Any(), gomock.Any(), data).
		Return("", errors.New("test error"))

	// Вызов функции
	HandleCreateShortURL(context.Background(), w, req, data, baseHost, zapSugar, shortURLService)

	// Проверка кода ответа
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
