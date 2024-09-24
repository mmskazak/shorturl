package mocks

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

// MockBatchResults - это мок для интерфейса BatchResults.
type MockBatchResults struct {
	pgx.BatchResults
	mock.Mock
}

// QueryRow - запрос строки из mock базы данных.
func (m *MockBatchResults) QueryRow() pgx.Row {
	args := m.Called()
	row, ok := args.Get(0).(pgx.Row)
	if !ok {
		return nil
	}
	return row
}

// Close - закрытие соединения с mock базой данных.
func (m *MockBatchResults) Close() error {
	args := m.Called()
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("failed to close mock batch results: %w", err)
	}
	return nil
}
