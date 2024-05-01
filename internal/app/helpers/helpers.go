package helpers

import (
	"math/rand"
	"time"
)

// GenerateShortURL генерирует случайный строковый идентификатор заданной длины.
func GenerateShortURL(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	// Создаем новый генератор случайных чисел
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	for i := range b {
		b[i] = charset[rng.Intn(len(charset))] // генерация случайного индекса
	}
	return string(b)
}
