package inmemory

import (
	"sync"

	"go.uber.org/zap"
)

// URLRecord - структура для хранения URL с дополнительной информацией.
type URLRecord struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"` // Оригинальный URL
	UserID      string `json:"user_id"`      // Идентификатор пользователя
	Deleted     bool   `json:"deleted"`      // Флаг, указывающий на удаление URL
}

// InMemory - структура для работы с хранилищем в памяти.
type InMemory struct {
	Mu        *sync.Mutex
	Data      map[string]URLRecord
	userIndex map[string][]string // для быстрого поиска URL по userID
	zapLog    *zap.SugaredLogger
}

// NewInMemory - конструктор для создания нового хранилища в памяти.
func NewInMemory(zapLog *zap.SugaredLogger) (*InMemory, error) {
	return &InMemory{
		Mu:        &sync.Mutex{},
		Data:      make(map[string]URLRecord),
		userIndex: make(map[string][]string),
		zapLog:    zapLog,
	}, nil
}

// Close - закрытие хранилища (заглушка для будущих изменений).
func (m *InMemory) Close() error {
	m.zapLog.Debugln("InMemory storage closed (nothing to close currently)")
	return nil
}

// NumberOfEntries - количество записей.
func (m *InMemory) NumberOfEntries() int {
	return len(m.Data)
}
