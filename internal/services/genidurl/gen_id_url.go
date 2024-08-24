package genidurl

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenID предоставляет методы для генерации уникальных идентификаторов.
type GenID struct{}

// Generate создает случайный идентификатор заданной длины.
func (s *GenID) Generate() (string, error) {
	const length = 8                                                                 // Длина генерируемого идентификатора
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Набор символов для генерации

	b := make([]byte, length) // Создаем срез для хранения идентификатора

	for i := range b {
		// Генерируем случайный индекс для выбора символа из charset
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("ошибка генерации случайного индекса: %w", err)
		}
		b[i] = charset[randomIndex.Int64()] // Присваиваем символ из charset в срез
	}
	return string(b), nil // Преобразуем срез в строку и возвращаем
}

// NewGenIDService создает новый экземпляр генератора идентификаторов.
func NewGenIDService() *GenID {
	return &GenID{}
}
