package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"

	"mmskazak/shorturl/internal/services/shorturlservice/mocks"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap/zaptest"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage(t *testing.T) {
	// Create a new chi router
	r := chi.NewRouter()

	// Define the route and bind the handler
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger, err := zap.NewDevelopment()
		require.NoError(t, err)
		zapLog := logger.Sugar()
		MainPage(w, r, zapLog)
	})

	// Создаем фейковый HTTP запрос
	req, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	require.NoError(t, err)

	// Create a new response recorder to capture the response from the handler
	w := httptest.NewRecorder()

	// Call the handler function with the test request and response recorder
	r.ServeHTTP(w, req)

	// Check that the response status code is 201 Created
	require.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, w.Body.String(), "Сервис сокращения URL")
}

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

	HandleCreateShortURL(ctxBg, w, r, data, baseHost, zapSugar)

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

	// Вызов функции
	HandleCreateShortURL(context.Background(), w, req, data, baseHost, zapSugar)

	// Проверка кода ответа
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreateShortURL(t *testing.T) {
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
	body := strings.NewReader("http://yandex.ru")
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockISetShortURL(ctrl)
	data.EXPECT().SetShortURL(
		gomock.Any(),
		gomock.Any(),
		"http://yandex.ru",
		"11111",
		false)

	// Базовый хост
	baseHost := "http://localhost"

	// Вызов функции
	HandleCreateShortURL(context.Background(), w, req, data, baseHost, zapSugar)

	// Проверка кода ответа
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, w.Body.String())
}
