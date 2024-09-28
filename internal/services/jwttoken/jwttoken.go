package jwttoken

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"mmskazak/shorturl/internal/services/jwtbuilder"
)

// CreateNewJWTToken - создает новый JWT токен с новым id пользователя
func CreateNewJWTToken(userID, secretKey string) (string, error) {
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
		return "", fmt.Errorf("error func CreateJWTToken on "+
			"jwt.Create(header, payloadStruct, secretKey): %w", err)
	}

	return token, nil
}

// GetSignedPayloadJWT извлекает и проверяет подписанную полезную нагрузку JWT из cookie.
func GetSignedPayloadJWT(jwt string, secretKey string) (string, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 { //nolint:gomnd // 3 части JWT токена
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

// verifyHMAC проверяет, соответствует ли предоставленная подпись ожидаемому значению HMAC.
func verifyHMAC(value, signature, key string) bool {
	expectedSignature := jwtbuilder.GenerateHMAC(value, key)
	return compareHMAC(expectedSignature, signature)
}

// compareHMAC сравнивает два HMAC значения, возвращая true, если они идентичны.
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

func GetUserIDFromJWT(jwt string, secretKey string) (string, error) {
	payloadString, err := GetSignedPayloadJWT(jwt, secretKey)
	if err != nil {
		return "", fmt.Errorf("error GetUserIDFromJWT: %w", err)
	}
	payloadStruct := jwtbuilder.PayloadJWT{}
	err = json.Unmarshal([]byte(payloadString), &payloadStruct)
	if err != nil {
		return "", fmt.Errorf("error Unmarshal payload: %w", err)
	}

	return payloadStruct.UserID, nil
}
