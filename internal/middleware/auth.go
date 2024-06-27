package middleware

import (
	"context"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
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
	authorizationCookieName = "authorization"
)

func compareHMAC(sig1, sig2 string) bool {
	decodedSig1, err := base64.RawURLEncoding.DecodeString(sig1)
	if err != nil {
		return false
	}

	decodedSig2, err := base64.RawURLEncoding.DecodeString(sig2)
	if err != nil {
		return false
	}

	return hmac.Equal(decodedSig1, decodedSig2)
}

func verifyHMAC(value, signature, key string) bool {
	expectedSignature := jwtbuilder.GenerateHMAC(value, key)
	return compareHMAC(expectedSignature, signature)
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
		var payloadStruct jwtbuilder.PayloadJWT
		payloadString, err := getSignedPayloadJWT(r, authorizationCookieName, secretKey)
		if err != nil {
			// Логируем ошибку
			zapLog.Warnf("Failed to get signed payloadString of JWT: %v", err)

			// Создаем новый JWT токен
			userID := uuid.New().String()

			// Используем jwtbuilder для создания нового токена
			jwt := jwtbuilder.New()
			header := jwtbuilder.HeaderJWT{
				Alg: "HS256", // Укажите используемый вами алгоритм
				Typ: "JWT",
			}
			payloadStruct = jwtbuilder.PayloadJWT{
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
		} else {
			payloadStruct = jwtbuilder.PayloadJWT{}
			err = json.Unmarshal([]byte(payloadString), &payloadStruct)
			if err != nil {
				zapLog.Errorf("error unmarshalling payloadString: %v", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		zapLog.Infof("Payload to context: %s", payloadStruct)
		// Добавляем payloadString в контекст
		ctx := context.WithValue(r.Context(), ctxkeys.PayLoad, payloadStruct)
		r = r.WithContext(ctx)

		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}
