package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/contracts/mocks"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"mmskazak/shorturl/internal/services/shorturlservice"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
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

func TestPingPostgreSQL(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	pinger := mocks.NewMockPinger(ctrl)

	type args struct {
		w      http.ResponseWriter
		req    *http.Request
		data   contracts.Pinger
		zapLog *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test 1",
			args: args{
				w:      httptest.NewRecorder(),
				req:    httptest.NewRequest(http.MethodPost, "/", http.NoBody),
				data:   pinger,
				zapLog: zap.NewNop().Sugar(),
			},
		},
	}
	for _, tt := range tests {
		if tt.name == "test 1" {
			pinger.EXPECT().Ping(ctx)
		}
		t.Run(tt.name, func(t *testing.T) {
			PingPostgreSQL(ctx, tt.args.w, tt.args.req, tt.args.data, tt.args.zapLog)
		})
	}
}
