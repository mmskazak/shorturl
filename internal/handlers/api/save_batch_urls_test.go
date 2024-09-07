package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	storageErrors "mmskazak/shorturl/internal/storage/errors"

	"mmskazak/shorturl/internal/contracts/mocks"
	"mmskazak/shorturl/internal/models"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/jwtbuilder"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
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

func TestSaveShortenURLsBatch_ErrStatusUnauthorized(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание контекста с PayloadJWT
	ctx := context.Background()

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

	ctxBg := context.Background()

	// Настройка роутера chi для обработки запроса
	router := chi.NewRouter()
	router.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		SaveShortenURLsBatch(ctxBg, w, r, store, baseHost, zapSugar)
	})

	// Выполнение запроса через роутер
	router.ServeHTTP(w, req)

	// Проверка кода ответа
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSaveShortenURLsBatch_ErrSaveBatch(t *testing.T) {
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
	store.EXPECT().SaveBatch(ctxBg, incomingOutput, baseHost, "11111", generator).
		Return(nil, errors.New("test error"))

	// Настройка роутера chi для обработки запроса
	router := chi.NewRouter()
	router.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		SaveShortenURLsBatch(ctxBg, w, r, store, baseHost, zapSugar)
	})

	// Выполнение запроса через роутер
	router.ServeHTTP(w, req)

	// Проверка кода ответа
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSaveShortenURLsBatch_ErrUniqueViolation(t *testing.T) {
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
	store.EXPECT().SaveBatch(ctxBg, incomingOutput, baseHost, "11111", generator).
		Return(nil, storageErrors.ErrUniqueViolation)

	// Настройка роутера chi для обработки запроса
	router := chi.NewRouter()
	router.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		SaveShortenURLsBatch(ctxBg, w, r, store, baseHost, zapSugar)
	})

	// Выполнение запроса через роутер
	router.ServeHTTP(w, req)

	// Проверка кода ответа
	assert.Equal(t, http.StatusConflict, w.Code)
}
