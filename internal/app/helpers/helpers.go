package helpers

import (
	"crypto/rand"
	"math/big"
)

// GenerateShortURL генерирует случайный строковый идентификатор заданной длины.
func GenerateShortURL(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	for i := range b {
		// Генерируем случайный индекс для выбора символа из charset
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err) // Обрабатываем ошибку, если не удается сгенерировать случайное число
		}
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}
