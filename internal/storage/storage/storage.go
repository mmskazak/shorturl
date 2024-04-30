package storage

import (
	"errors"
)

var ErrNotFound = errors.New("key not found")

type Storage interface {
	GetShortURL(id string) string            //получение короткого URL
	SetShortURL(id string, targetURL string) //установить короткий URL
}
