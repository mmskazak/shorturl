package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"mmskazak/shorturl/internal/services/jwttoken"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

const (
	authorizationCookieName = "authorization"
)

// setSignedCookie устанавливает подписанный JWT в виде cookie в ответе HTTP.
func setSignedCookie(w http.ResponseWriter, name string, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})
}

// AuthMiddleware создает middleware для аутентификации запросов на основе JWT cookie.
// Если JWT отсутствует или недействителен, создается новый JWT токен и устанавливается в cookie.
func AuthMiddleware(next http.Handler, cfg *config.Config, zapLog *zap.SugaredLogger) http.Handler {
	secretKey := cfg.SecretKey
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payloadStruct jwtbuilder.PayloadJWT
		payloadString, err := getSignedJWTPayloadFromCookie(r, secretKey)
		if err != nil {
			// Логируем ошибку при получении или проверке JWT
			zapLog.Warnf("Failed to get signed payloadString of JWT: %v", err)

			// Создаем новый userID
			userID := uuid.New().String()

			token, err := jwttoken.CreateNewJWTToken(userID, secretKey)
			if err != nil {
				http.Error(w, "Failed to create new authorization token", http.StatusInternalServerError)
				return
			}
			payloadString, err = jwttoken.GetSignedPayloadJWT(token, secretKey)
			if err != nil {
				http.Error(w, "Failed to create new authorization token", http.StatusInternalServerError)
				return
			}

			// Устанавливаем в куку новый jwt
			setSignedCookie(w, authorizationCookieName, token)
		}

		err = json.Unmarshal([]byte(payloadString), &payloadStruct)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		zapLog.Infof("Payload isset: %s", payloadStruct)

		zapLog.Infof("Payload to context before install: %s", payloadStruct)
		// Добавляем payloadString в контекст
		ctx := context.WithValue(r.Context(), ctxkeys.PayLoad, payloadStruct)
		r = r.WithContext(ctx)

		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}

func getSignedJWTPayloadFromCookie(r *http.Request, secretKey string) (string, error) {
	cookie, err := r.Cookie(authorizationCookieName)
	if err != nil {
		return "", fmt.Errorf("no authorization cookie: %w", err)
	}
	jwt := cookie.Value
	payloadString, err := jwttoken.GetSignedPayloadJWT(jwt, secretKey)
	if err != nil {
		return "", fmt.Errorf("error getting signed payloadString of JWT: %w", err)
	}

	return payloadString, nil
}
