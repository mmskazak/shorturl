package errors

import "errors"

var (
	ErrKeyAlreadyExists = errors.New("key already exists")
	ErrNotFound         = errors.New("key not found")
)
