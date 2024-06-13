package errors

import (
	"errors"
	"fmt"
)

var (
	ErrOriginalURLAlreadyExists = errors.New("error original url already exists")
	ErrKeyAlreadyExists         = errors.New("error key already exists")
	ErrNotFound                 = errors.New("error key not found")
	ErrUniqueViolation          = errors.New("error database unique violation")
)

type ConflictError struct {
	Err      error
	ShortURL string
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("%v: %v", ErrOriginalURLAlreadyExists, e.ShortURL)
}
