package mocks

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

// MockDatabase - это структура для мока интерфейса Database
type MockDatabase struct {
	pgxpool.Pool
	mock.Mock
}

// Exec - мок для метода Exec
func (m *MockDatabase) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	a := m.Called(ctx, sql, args)
	return a.Get(0).(pgconn.CommandTag), a.Error(1)
}

// QueryRow - мок для метода QueryRow
func (m *MockDatabase) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	a := m.Called(ctx, sql, args)
	return a.Get(0).(pgx.Row)
}

// Ping - мок для метода Ping
func (m *MockDatabase) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Begin - мок для метода Begin
func (m *MockDatabase) Begin(ctx context.Context) (pgx.Tx, error) {
	a := m.Called(ctx)
	return a.Get(0).(pgx.Tx), a.Error(1)
}
