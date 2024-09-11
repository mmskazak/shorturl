package mocks

import (
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

// MockBatchResults - это мок для интерфейса BatchResults
type MockBatchResults struct {
	pgx.BatchResults
	mock.Mock
}

func (m *MockBatchResults) QueryRow() pgx.Row {
	args := m.Called()
	return args.Get(0).(pgx.Row)
}

func (m *MockBatchResults) Close() error {
	args := m.Called()
	return args.Error(0)
}
