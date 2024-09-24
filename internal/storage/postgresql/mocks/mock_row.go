package mocks

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

// MockRow — структура, представляющая мок для интерфейса Row.
type MockRow struct {
	pgx.Row
	mock.Mock
}

// Scan Реализация метода Scan для mock.
func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error MockRow func Scan: %w", err)
	}
	return nil
}
