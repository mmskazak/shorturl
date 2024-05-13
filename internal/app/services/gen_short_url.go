package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type GenShortURL struct {
}

func (s *GenShortURL) Generate(length int) (string, error) {
	const minLengthShortURL = 4
	if length < minLengthShortURL {
		return "", errors.New("length short URl too small")
	}

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

func NewGenIDService() *GenShortURL {
	return &GenShortURL{}
}
