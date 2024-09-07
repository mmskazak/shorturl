package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mmskazak/shorturl/internal/contracts/mocks"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestDeleteUserURLs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "55555"})

	zapSugar := zap.NewNop().Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()
	// Создание HTTP-запроса с телом запроса

	body := strings.NewReader(`
[
  "BsfN139Y",
  "FLpqbWug"
]
`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockIDeleteUserURLs(ctrl)
	expectParam := []string{"BsfN139Y", "FLpqbWug"}

	data.EXPECT().DeleteURLs(context.Background(), expectParam).Return(nil)

	// Вызов функции
	DeleteUserURLs(context.Background(), w, req, data, zapSugar)

	// Проверка кода ответа
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestDeleteUserURLs_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT
	ctx := context.WithValue(context.Background(), ctxkeys.PayLoad, jwtbuilder.PayloadJWT{UserID: "55555"})

	zapSugar := zap.NewNop().Sugar()

	// Создание HTTP-запроса и ResponseRecorder
	w := httptest.NewRecorder()
	// Создание HTTP-запроса с телом запроса

	body := strings.NewReader(`
[
  "BsfN139Y",
  "FLpqbWug"
]
`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req = req.WithContext(ctx)

	// Создание мока для ISetShortURL
	data := mocks.NewMockIDeleteUserURLs(ctrl)
	expectParam := []string{"BsfN139Y", "FLpqbWug"}

	data.EXPECT().DeleteURLs(context.Background(), expectParam).
		Return(errors.New("test errors"))

	// Вызов функции
	DeleteUserURLs(context.Background(), w, req, data, zapSugar)

	// Проверка кода ответа
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
