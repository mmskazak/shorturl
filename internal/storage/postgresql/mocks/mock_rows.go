package mocks

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

// MockRows представляет мок для интерфейса Rows.
type MockRows struct {
	pgx.Rows
	mock.Mock
}

// Close Реализация метода.
func (m *MockRows) Close() {
	m.Called()
}

// Err Реализация метода.
func (m *MockRows) Err() error {
	args := m.Called()
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error MockRows func Err: %w", err)
	}
	return nil
}

// Next Реализация метода.
func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

// Scan Реализация метода.
func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error MockRows func Scan: %w", err)
	}
	return nil
}
