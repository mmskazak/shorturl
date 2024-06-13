package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	cookieName = "user_id"
	secretKey  = "supersecretkey"
)

func generateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func verifyHMAC(data, signature, key string) bool {
	expectedSignature := generateHMAC(data, key)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// setSignedCookie sets a cookie with HMAC signature.
func setSignedCookie(w http.ResponseWriter, name, value, key string) {
	signature := generateHMAC(value, key)
	cookieValue := fmt.Sprintf("%s.%s", value, signature)
	cookie := &http.Cookie{
		Name:     name,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
}

func getSignedCookie(r *http.Request, name, key string) (string, bool) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", false
	}

	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return "", false
	}

	value, signature := parts[0], parts[1]
	if verifyHMAC(value, signature, key) {
		return value, true
	}
	return "", false
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, valid := getSignedCookie(r, cookieName, secretKey)
		log.Printf("User ID: %s Valid: %t", userID, valid)
		if !valid {
			userID = uuid.New().String()
			setSignedCookie(w, cookieName, userID, secretKey)
		}
		next.ServeHTTP(w, r)
	})
}
