package jwttoken

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"fmt"
	jwtGolang "github.com/golang-jwt/jwt"
	"mmskazak/shorturl/internal/services/jwtbuilder"
)

// CreateNewJWTToken - создает новый JWT токен с новым id пользователя.
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
func GetSignedPayloadJWT(tokenString string, secretKey string) (string, error) {
	token, err := jwtGolang.ParseWithClaims(tokenString, &jwtbuilder.PayloadJWT{}, func(token *jwtGolang.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtGolang.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("error GetSignedPayloadJWT: %w", err)
	}

	data, err := json.Marshal(token.Claims)
	if err != nil {
		return "", fmt.Errorf("error marshalling token claims: %w", err)
	}

	return string(data), nil
}

// verifyHMAC проверяет, соответствует ли предоставленная подпись ожидаемому значению HMAC.
func verifyHMAC(value, signature, key string) bool {
	expectedSignature := jwtbuilder.GenerateHMAC(value, key)
	return CompareHMAC(expectedSignature, signature)
}

// CompareHMAC сравнивает два HMAC значения, возвращая true, если они идентичны.
func CompareHMAC(sig1, sig2 string) bool {
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
