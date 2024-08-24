package jwtbuilder

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// HeaderJWT представляет заголовок JWT, который содержит информацию о типе токена и алгоритме подписи.
type HeaderJWT struct {
	Alg string `json:"alg"` // Алгоритм подписи (например, HS256)
	Typ string `json:"typ"` // Тип токена (например, JWT)
}

// PayloadJWT представляет полезную нагрузку JWT, которая содержит данные пользователя.
type PayloadJWT struct {
	UserID string `json:"user_id"` // Идентификатор пользователя
}

// JWTBuilder представляет собой конструктор для создания JWT.
type JWTBuilder struct{}

// New создает новый экземпляр JWTBuilder.
func New() JWTBuilder {
	return JWTBuilder{}
}

// Create формирует JWT, используя заголовок, полезную нагрузку и секретный ключ.
func (j *JWTBuilder) Create(header HeaderJWT, payload PayloadJWT, secret string) (string, error) {
	// Сериализация структур в JSON
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Кодирование JSON в Base64 URL
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Создание подписи HMAC
	signature := GenerateHMAC(fmt.Sprintf("%s.%s", encodedHeader, encodedPayload), secret)

	// Формирование окончательного токена
	token := fmt.Sprintf("%s.%s.%s", encodedHeader, encodedPayload, signature)
	return token, nil
}

// GenerateHMAC создает HMAC-подпись для данных, используя секретный ключ.
func GenerateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
