package web

import (
	"bytes"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"mmskazak/shorturl/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage(t *testing.T) {
	// Create a new chi router
	r := chi.NewRouter()

	// Define the route and bind the handler
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r)
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

func TestCreateShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок-объект для Storage
	mockStorage := mocks.NewMockStorage(ctrl)

	// Определяем, что мы ожидаем вызов SetShortURL с аргументами и возвращаем nil ошибку
	mockStorage.EXPECT().SetShortURL(gomock.Any(), "http://ya.ru").Return(nil)

	// Создаем хендлер с мок-объектом
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCreateShortURL(w, r, mockStorage, "http://localhost")
	})

	// Создаем новый запрос с методом POST и телом, содержащим оригинальный URL
	reqBody := bytes.NewBufferString("http://ya.ru")
	req := httptest.NewRequest(http.MethodPost, "/", reqBody)

	// Создаем новый response recorder для захвата ответа от хендлера
	w := httptest.NewRecorder()

	// Вызываем хендлер с тестовым запросом и response recorder
	handler.ServeHTTP(w, req)

	// Проверяем, что статус ответа 201 Created
	require.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

func TestHandleRedirect(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		expectedCode int
		mockSetup    func(ms *mocks.MockStorage)
	}{
		{
			name:         "NotFound",
			path:         "/x0x0x0x0",
			expectedCode: http.StatusNotFound,
			mockSetup: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetShortURL("x0x0x0x0").Return("", storageErrors.ErrNotFound)
			},
		},
		{
			name:         "BadRequest",
			path:         "/x0x0",
			expectedCode: http.StatusNotFound,
			mockSetup: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetShortURL("x0x0").Return("", storageErrors.ErrNotFound)
			},
		},
		{
			name:         "Redirect",
			path:         "/vAlIdIds",
			expectedCode: http.StatusTemporaryRedirect,
			mockSetup: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetShortURL("vAlIdIds").Return("http://ya.ru", nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаем мок-объект для Storage
			ms := mocks.NewMockStorage(ctrl)

			// Устанавливаем ожидания для текущего тест-кейса
			tc.mockSetup(ms)

			handleRedirectHandler := func(w http.ResponseWriter, r *http.Request) {
				HandleRedirect(w, r, ms)
			}

			// Создаем новый роутер и регистрируем хендлер
			r := chi.NewRouter()
			r.Get("/{id}", handleRedirectHandler)

			// Создаем новый запрос с текущим path
			req, err := http.NewRequest(http.MethodGet, tc.path, http.NoBody)
			require.NoError(t, err)

			// Создаем новый response recorder для захвата ответа от хендлера
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			// Проверяем код ответа
			assert.Equal(t, tc.expectedCode, rr.Code)
		})
	}
}
