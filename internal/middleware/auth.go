package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"

	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	cookieName = "user_id"
	secretKey  = "supersecretkey"
)

// Определяем тип для ключа контекста.
type contextKey string

// Постоянный ключ для идентификатора пользователя.
const keyUserID contextKey = "userID"

func generateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil))
}

func verifyHMAC(value, signature, key string) bool {
	expectedSignature := generateHMAC(value, key)
	log.Println("Value for HMAC generation:", value)
	log.Println("Expected Signature:", expectedSignature)
	log.Println("Provided Signature:", signature)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// setSignedCookie sets a cookie with HMAC signature.
func setSignedCookie(w http.ResponseWriter, name, value, key string) {
	signature := generateHMAC(value, key)
	cookieValue := fmt.Sprintf("%s@%s", value, signature)
	encodedValue := url.QueryEscape(cookieValue)
	cookie := &http.Cookie{
		Name:     name,
		Value:    encodedValue,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour), //nolint:gomnd //24 часа
	}
	http.SetCookie(w, cookie)
}

func getSignedCookie(r *http.Request, name, key string) (string, bool) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", false
	}

	decodedCookie, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		log.Printf("Error decoding cookie value: %v", err)
		return "", false
	}

	parts := strings.Split(decodedCookie, "@")
	if len(parts) != 2 { //nolint:gomnd // uuid и подпись
		log.Println("Invalid cookie format")
		return "", false
	}

	decodedValue := parts[0]
	decodedSignature := parts[1]

	if verifyHMAC(decodedValue, decodedSignature, key) {
		return decodedValue, true
	}

	return "", false
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, valid := getSignedCookie(r, cookieName, secretKey)
		log.Printf("User ID: %s Valid: %t", userID, valid)

		// Если кука недействительна или отсутствует
		if !valid {
			userID = uuid.New().String()
			setSignedCookie(w, cookieName, userID, secretKey)
			// Получаем URI текущего запроса
			uri := r.URL.Path
			log.Printf("Request URI: %s", uri)
			if uri == "/api/user/urls" {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
		}

		ctx := context.WithValue(r.Context(), keyUserID, userID)
		r = r.WithContext(ctx)

		// Если кука действительна, продолжаем выполнение следующего обработчика
		next.ServeHTTP(w, r)
	})
}
