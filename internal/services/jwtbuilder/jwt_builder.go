package jwtbuilder

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type JWTBuilder struct{}

func (j *JWTBuilder) Create(header, playLoad, secret string) (string, error) {

	return ""
}

func generateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil))
}
