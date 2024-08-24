package api

import (
	"context"
	"mmskazak/shorturl/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/handlers/api/mocks"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/jwtbuilder"
)

func TestSaveShortenURLsBatch(t *testing.T) {
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
	body := strings.NewReader(`
[
  {
    "correlation_id": "123",
    "original_url": "https://example.com/long-url-00012"
  },
  {
    "correlation_id": "456",
    "original_url": "https://example.com/long-url-00013"
  }
]
`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", body)
	req = req.WithContext(ctx)

	// Базовый хост
	baseHost := "http://localhost"

	store := mocks.NewMockISaveBatch(ctrl)
	generator := genidurl.NewGenIDService()

	expectedOutput := []models.Output{
		{
			CorrelationID: "123",
			ShortURL:      "http://localhost:8080/nQm6WEim",
		},
		{
			CorrelationID: "456",
			ShortURL:      "http://localhost:8080/2RWV924w",
		},
	}

	incomingOutput := []models.Incoming{
		{
			CorrelationID: "123",
			OriginalURL:   "https://example.com/long-url-00012",
		},
		{
			CorrelationID: "456",
			OriginalURL:   "https://example.com/long-url-00013",
		},
	}
	ctxBg := context.Background()
	store.EXPECT().SaveBatch(ctxBg, incomingOutput, baseHost, "11111", generator).Return(expectedOutput, nil)

	// Настройка роутера chi для обработки запроса
	router := chi.NewRouter()
	router.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		SaveShortenURLsBatch(ctxBg, w, r, store, baseHost, zapSugar)
	})

	// Выполнение запроса через роутер
	router.ServeHTTP(w, req)

	// Проверка кода ответа
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "[{\"correlation_id\":\"123\",\"short_url\":\"http://localhost:8080/nQm6WEim\"},"+
		"{\"correlation_id\":\"456\",\"short_url\":\"http://localhost:8080/2RWV924w\"}]",
		w.Body.String())
}
