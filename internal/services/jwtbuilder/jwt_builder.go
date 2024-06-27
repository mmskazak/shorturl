package jwtbuilder

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type HeaderJWT struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type PayloadJWT struct {
	UserID string `json:"user_id"`
}

type JWTBuilder struct{}

func New() JWTBuilder {
	return JWTBuilder{}
}

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
	signature := generateHMAC(fmt.Sprintf("%s.%s", encodedHeader, encodedPayload), secret)

	// Формирование окончательного токена
	token := fmt.Sprintf("%s.%s.%s", encodedHeader, encodedPayload, signature)
	return token, nil
}

func generateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
