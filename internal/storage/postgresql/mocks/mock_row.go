package mocks

import (
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

// MockRow — структура, представляющая мок для интерфейса Row
type MockRow struct {
	pgx.Row
	mock.Mock
}

// Реализация метода Scan для мока.
func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}
