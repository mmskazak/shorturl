package storage

import (
	"errors"
)

var ErrNotFound = errors.New("key not found")

type Repositories interface {
	GetShortURL(id string) string            // получение короткого URL
	SetShortURL(id string, targetURL string) // установить короткий URL
}
