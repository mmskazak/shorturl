package mocks

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockBatchResults - это мок для интерфейса BatchResults
type MockBatchResults struct {
	mock.Mock
}

func (m *MockBatchResults) Exec() (pgconn.CommandTag, error) {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *MockBatchResults) Query() (pgx.Rows, error) {
	args := m.Called()
	return args.Get(0).(pgx.Rows), args.Error(1)
}

func (m *MockBatchResults) QueryRow() pgx.Row {
	args := m.Called()
	return args.Get(0).(pgx.Row)
}

func (m *MockBatchResults) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockBatchResults) CloseWithError(ctx context.Context, err error) error {
	args := m.Called(ctx, err)
	return args.Error(0)
}
