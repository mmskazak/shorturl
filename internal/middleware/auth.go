package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/ctxkeys"
	"net/url"

	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	authorizationCookie = "Authorization"
)

type playLoadStruct struct {
	UserID string `json:"user_id"`
}

func generateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil))
}

func verifyHMAC(value, signature, key string) bool {
	expectedSignature := generateHMAC(value, key)
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

func getSignedPlayLoadJWT(r *http.Request, name, secretKey string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("error get signed cookie jwt: %w", err)
	}

	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid sruct jwt")
	}
	headerStr := parts[0]
	playLoadStr := parts[1]
	signatureStr := parts[2]

	if verifyHMAC(headerStr+"."+playLoadStr, signatureStr, secretKey) {
		playLoadStrDecoded, err := base64.RawURLEncoding.DecodeString(playLoadStr)
		if err != nil {
			return "", fmt.Errorf("error decode playLoadStr: %w", err)
		}
		return string(playLoadStrDecoded), nil
	}

	return "", fmt.Errorf("invalid verify hmac signature")
}

func AuthMiddleware(next http.Handler, cfg *config.Config, zapLog zap.SugaredLogger) http.Handler {
	secretKey := cfg.SecretKey
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playLoad, err := getSignedPlayLoadJWT(r, authorizationCookie, secretKey)
		if err != nil {
			// Если у нас ошибка тогда мы должны выдать JWT токен
		}

		// Если мы получили структуру то в ней долджен быть userID
		// в том стлуче если его нет то у нас тогда не стработает json.Unmarshal или что?
		pl := playLoadStruct{}
		err = json.Unmarshal([]byte(playLoad), &pl)
		if err != nil {
			zapLog.Errorf("error unmarshal playLoad struct: %v", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		userID := pl.UserID

		if err != nil {
			userID = uuid.New().String()
			setSignedCookie(w, authorizationCookie, userID, secretKey)
			// Получаем URI текущего запроса
			uri := r.URL.Path
			log.Printf("Request URI: %s", uri)
			if uri == "/api/user/urls" {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
		}

		ctx := context.WithValue(r.Context(), ctxkeys.KeyUserID, userID)
		r = r.WithContext(ctx)

		// Если кука действительна, продолжаем выполнение следующего обработчика
		next.ServeHTTP(w, r)
	})
}
