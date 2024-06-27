package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	authorizationCookieName = "Authorization"
)

func generateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil))
}

func verifyHMAC(value, signature, key string) bool {
	expectedSignature := generateHMAC(value, key)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func setSignedCookie(w http.ResponseWriter, name string, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})
}

func getSignedPayloadJWT(r *http.Request, name, secretKey string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("error get signed cookie jwt: %w", err)
	}

	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 3 { //nolint:gomnd //3 части jwt токена
		return "", errors.New("invalid structure jwt")
	}
	headerStr, payloadStr, signatureStr := parts[0], parts[1], parts[2]

	// Проверка HMAC подписи
	if !verifyHMAC(headerStr+"."+payloadStr, signatureStr, secretKey) {
		return "", errors.New("invalid HMAC signature verification")
	}

	// Декодирование полезной нагрузки из Base64 URL
	decodedPayload, err := base64.RawURLEncoding.DecodeString(payloadStr)
	if err != nil {
		return "", fmt.Errorf("error decoding payload: %w", err)
	}

	return string(decodedPayload), nil
}

func AuthMiddleware(next http.Handler, cfg *config.Config, zapLog *zap.SugaredLogger) http.Handler {
	secretKey := cfg.SecretKey
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := getSignedPayloadJWT(r, authorizationCookieName, secretKey)
		if err != nil {
			// Логируем ошибку
			zapLog.Warnf("Failed to get signed payload JWT: %v", err)

			// Создаем новый JWT токен
			userID := uuid.New().String()

			// Используем jwtbuilder для создания нового токена
			jwt := jwtbuilder.New()
			header := jwtbuilder.HeaderJWT{
				Alg: "HS256", // Укажите используемый вами алгоритм
				Typ: "JWT",
			}
			payloadStruct := jwtbuilder.PayloadJWT{
				UserID: userID,
			}

			token, err := jwt.Create(header, payloadStruct, secretKey)
			if err != nil {
				zapLog.Errorf("Failed to create JWT: %v", err)
				http.Error(w, "Failed to create authorization token", http.StatusInternalServerError)
				return
			}

			// Устанавливаем JWT токен в куки
			setSignedCookie(w, authorizationCookieName, token)

			zapLog.Infof("Issued new JWT for user: %s", userID)
		}

		// Добавляем payload в контекст
		ctx := context.WithValue(r.Context(), ctxkeys.PayLoad, payload)
		r = r.WithContext(ctx)

		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}
