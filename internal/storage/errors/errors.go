package errors

import "errors"

var (
	ErrOriginalURLAlreadyExists = errors.New("original url already exists")
	ErrKeyAlreadyExists         = errors.New("key already exists")
	ErrNotFound                 = errors.New("key not found")
)
