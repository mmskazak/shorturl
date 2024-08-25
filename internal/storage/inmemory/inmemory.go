package inmemory

import (
	"sync"

	"mmskazak/shorturl/internal/models"

	"go.uber.org/zap"
)

// InMemory - структура для работы с хранилищем в памяти.
type InMemory struct {
	mu        *sync.Mutex
	data      map[string]models.URLRecord
	userIndex map[string][]string // для быстрого поиска URL по userID
	zapLog    *zap.SugaredLogger
}

// NewInMemory - конструктор для создания нового хранилища в памяти.
func NewInMemory(zapLog *zap.SugaredLogger) (*InMemory, error) {
	return &InMemory{
		mu:        &sync.Mutex{},
		data:      make(map[string]models.URLRecord),
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
	return len(m.data)
}

// GetCopyData - копирует данные из памяти для того чтобы сохранить в файл.
func (m *InMemory) GetCopyData() map[string]models.URLRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Создаем новую карту данных
	copyData := make(map[string]models.URLRecord)

	// Копируем элементы из m.data в copyData
	for key, value := range m.data {
		copyData[key] = value
	}

	// Возвращаем копию данных
	return copyData
}
