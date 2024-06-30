package genidurl

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type GenID struct{}

func (s *GenID) Generate() (string, error) {
	length := 8
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	for i := range b {
		// Генерируем случайный индекс для выбора символа из charset
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("генерирование случайного индекса завершилось ошибкой %w", err)
		}
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b), nil
}

func NewGenIDService() *GenID {
	return &GenID{}
}
