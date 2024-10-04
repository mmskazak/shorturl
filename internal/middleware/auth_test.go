package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mmskazak/shorturl/internal/services/jwttoken"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func Test_compareHMAC(t *testing.T) {
	type args struct {
		sig1 string
		sig2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "identical HMACs",
			args: args{
				sig1: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue")),
				sig2: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue")),
			},
			want: true,
		},
		{
			name: "different HMACs",
			args: args{
				sig1: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue1")),
				sig2: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue2")),
			},
			want: false,
		},
		{
			name: "invalid base64 encoding for sig1",
			args: args{
				sig1: "invalidBase64",
				sig2: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue")),
			},
			want: false,
		},
		{
			name: "invalid base64 encoding for sig2",
			args: args{
				sig1: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue")),
				sig2: "invalidBase64",
			},
			want: false,
		},
		{
			name: "both HMACs are invalid",
			args: args{
				sig1: "invalidBase64",
				sig2: "invalidBase64",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, jwttoken.CompareHMAC(tt.args.sig1, tt.args.sig2),
				"compareHMAC(%v, %v)", tt.args.sig1, tt.args.sig2)
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "test_secret_key",
	}

	// Создаем фейковый логгер для тестов
	zapLog := zap.NewNop().Sugar()

	// Создаем HTTP тестовый сервер с применением middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем payloadString из контекста
		payload, ok := r.Context().Value(ctxkeys.PayLoad).(jwtbuilder.PayloadJWT)
		if !ok {
			// Если userID не найден или неверного типа, возвращаем ошибку
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// Сначала кодируем структуру в JSON
		jsonData, err := json.Marshal(payload)
		require.NoError(t, err)

		// Устанавливаем заголовок Content-Type для ответа
		w.Header().Set("Content-Type", "application/json")

		// Устанавливаем статус ответа
		w.WriteHeader(http.StatusOK)

		// Записываем JSON в тело ответа
		_, err = w.Write(jsonData)
		require.NoError(t, err)
	})

	rr := httptest.NewRecorder()
	middleware := AuthMiddleware(nextHandler, cfg, zapLog)

	// Тестовый случай 1: JWT отсутствует или недействителен
	req, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	middleware.ServeHTTP(rr, req)

	// Проверяем, что middleware создал новый JWT и установил его в куки
	assert.Equal(t, http.StatusOK, rr.Code)
	var responsePayload jwtbuilder.PayloadJWT
	err = json.Unmarshal(rr.Body.Bytes(), &responsePayload)
	require.NoError(t, err)
	assert.NotEmpty(t, responsePayload.UserID)

	// Тестовый случай 2: JWT уже существует и валиден
	// Создаем фейковый JWT и устанавливаем его в куки для теста
	jwt := jwtbuilder.New()
	header := jwtbuilder.HeaderJWT{
		Alg: "HS256",
		Typ: "JWT",
	}

	fakeUserID := "exampleUserID"

	payload := jwtbuilder.PayloadJWT{
		UserID: fakeUserID,
	}
	fakeToken, err := jwt.Create(header, payload, cfg.SecretKey)
	require.NoError(t, err)

	req.Header.Set("Cookie", authorizationCookieName+"="+fakeToken)
	rr = httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	// Проверяем, что middleware использовал существующий JWT из куки
	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.NewDecoder(rr.Body).Decode(&responsePayload)
	require.NoError(t, err)

	assert.Equal(t, fakeUserID, responsePayload.UserID) // заменить fakeUserID на фактический идентификатор пользователя
}
