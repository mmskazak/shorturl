package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mmskazak/shorturl/internal/contracts/mocks"
	"mmskazak/shorturl/internal/services/shorturlservice"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap/zaptest"

	"github.com/stretchr/testify/assert"
)

func TestHandleCreateShortURL_StatusUnauthorized(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxBg := context.Background()
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса с телом запроса
	body := strings.NewReader("http://yandex.ru")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", body)

	data := mocks.NewMockISetShortURL(ctrl)
	baseHost := "http://localhost"

	shortURLService := shorturlservice.NewShortURLService()

	HandleCreateShortURL(ctxBg, w, r, data, baseHost, zapSugar, shortURLService)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleCreateShortURL_EmptyBody(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockISetShortURL(ctrl)
	// Базовый хост
	baseHost := "http://localhost"

	shortURLService := shorturlservice.NewShortURLService()
	// Вызов функции
	HandleCreateShortURL(context.Background(), w, req, data, baseHost, zapSugar, shortURLService)

	// Проверка кода ответа
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreateShortURL_ErrUnmarshal(t *testing.T) {
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
	body := strings.NewReader(`{"url": "https://google.ru`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockISetShortURL(ctrl)

	// Базовый хост
	baseHost := "http://localhost"

	shortURLService := shorturlservice.NewShortURLService()

	// Вызов функции
	HandleCreateShortURL(context.Background(), w, req, data, baseHost, zapSugar, shortURLService)

	// Проверка кода ответа
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandleCreateShortURL_Success(t *testing.T) {
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
	data.EXPECT().SetShortURL(
		gomock.Any(),
		gomock.Any(),
		"https://google.ru",
		"11111",
		false)

	// Базовый хост
	baseHost := "http://localhost"

	shortURLService := shorturlservice.NewShortURLService()

	// Вызов функции
	HandleCreateShortURL(context.Background(), w, req, data, baseHost, zapSugar, shortURLService)

	// Проверка кода ответа
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

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
	fmt.Println(w.Body.String())
	assert.NotEmpty(t, w.Body.String())
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

func TestHandleCreateShortURL_ErrIoReadAll(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT(тут он не нужен)
	ctx := context.Background()

	// Создание логгера
	zapSugar := zaptest.NewLogger(t).Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()
	// Создание HTTP-запроса с телом запроса
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockISetShortURL(ctrl)

	// Базовый хост
	baseHost := "http://localhost"

	shortURLService := mocks.NewMockIShortURLService(ctrl)

	// Вызов функции
	HandleCreateShortURL(context.Background(), w, req, data, baseHost, zapSugar, shortURLService)

	// Проверка кода ответа
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
