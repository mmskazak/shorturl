package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrOriginalURLAlreadyExists указывает на ошибку, когда оригинальный URL уже существует в базе данных.
	// Эта ошибка возникает, если попытка добавить новый оригинальный URL завершилась неудачей из-за его дублирования.
	ErrOriginalURLAlreadyExists = errors.New("error original url already exists")

	// ErrKeyAlreadyExists указывает на ошибку, когда ключ (например,
	// уникальный идентификатор) уже существует в базе данных.
	// Эта ошибка может возникать, если попытка создать новый элемент с уже существующим ключом приводит к конфликту.
	ErrKeyAlreadyExists = errors.New("error key already exists")

	// ErrNotFound указывает на ошибку, когда ключ не найден в базе данных.
	// Эта ошибка возникает, если запрашиваемый элемент отсутствует в хранилище данных.
	ErrNotFound = errors.New("error key not found")

	// ErrUniqueViolation указывает на ошибку нарушения уникальности в базе данных.
	// Эта ошибка возникает, если в базу данных пытаются вставить дублирующие данные,
	// что нарушает ограничения уникальности.
	ErrUniqueViolation = errors.New("error database unique violation")

	// ErrShortURLsForUserNotFound указывает на ошибку,
	// когда короткие URL-адреса для указанного пользователя не найдены.
	// Эта ошибка может возникать при попытке получить короткие URL-адреса для пользователя,
	// у которого нет сохраненных URL.
	ErrShortURLsForUserNotFound = errors.New("short urls for user not found")

	// ErrDeletedShortURL указывает на ошибку, когда запрашиваемый короткий URL был удален.
	// Эта ошибка возникает, если попытка доступа к короткому URL не удалась, так как он был удален из системы.
	ErrDeletedShortURL = errors.New("deleted short url")
)

// ConflictError представляет ошибку конфликта, связанную с уникальностью короткого URL.
// Включает оригинальную ошибку и короткий URL, который вызвал конфликт.
type ConflictError struct {
	Err      error  // Оригинальная ошибка, связанная с конфликтом.
	ShortURL string // Короткий URL, который вызвал конфликт.
}

// Error возвращает строковое представление ошибки ConflictError.
// Включает оригинальную ошибку и короткий URL, что позволяет легче понять суть конфликта.
func (e ConflictError) Error() string {
	return fmt.Sprintf("%v: %v", ErrOriginalURLAlreadyExists, e.ShortURL)
}
