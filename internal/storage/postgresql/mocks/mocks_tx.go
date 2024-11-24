package mocks

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockTx - структура для мока интерфейса Tx.
type MockTx struct {
	pgx.Tx
	mock.Mock
}

// Commit Реализация метода.
func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error MockTx func Commit: %w", err)
	}
	return nil
}

// Rollback Реализация метода.
func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error MockTx func Rollback: %w", err)
	}
	return nil
}

// Exec Реализация метода.
func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	commandTag, _ := args.Get(0).(pgconn.CommandTag)
	if err := args.Error(1); err != nil {
		return commandTag, fmt.Errorf("error MockTx func Exec: %w", err)
	}

	return commandTag, nil
}

// Query Реализация метода.
func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	callArgs := m.Called(ctx, sql, args)
	rows, _ := callArgs.Get(0).(pgx.Rows)
	if err := callArgs.Error(1); err != nil {
		return rows, fmt.Errorf("error MockTx func Query: %w", err)
	}
	return rows, nil
}

// SendBatch Реализация метода.
func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, b)
	batchResults, _ := args.Get(0).(pgx.BatchResults)
	return batchResults
}

// QueryRow Реализация метода.
func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	// Получаем аргументы вызова
	call := m.Called(ctx, sql)
	// Извлекаем результат
	row := call.Get(0)
	// Возвращаем результат как pgx. Row
	if row == nil {
		return nil
	}
	r, ok := row.(pgx.Row)
	if !ok {
		return nil
	}
	return r
}
