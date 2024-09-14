package mocks

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

// MockDatabase - это структура для мока интерфейса Database.
type MockDatabase struct {
	pgxpool.Pool
	mock.Mock
}

// Exec - мок для метода Exec.
func (m *MockDatabase) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	a := m.Called(ctx, sql, args)
	cmdTag, ok := a.Get(0).(pgconn.CommandTag)
	if !ok {
		return pgconn.CommandTag{}, fmt.Errorf("failed to assert type for CommandTag: got %T", a.Get(0))
	}

	// Обработка и добавление контекста к ошибке
	if err := a.Error(1); err != nil {
		return cmdTag, fmt.Errorf("exec failed: %w", err)
	}

	return cmdTag, nil
}

// QueryRow - мок для метода QueryRow.
func (m *MockDatabase) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	a := m.Called(ctx, sql, args)
	// Приведение типа с проверкой
	row, ok := a.Get(0).(pgx.Row)
	if !ok {
		return nil
	}

	return row
}

// Ping - мок для метода Ping.
func (m *MockDatabase) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	err := args.Error(0)
	if err != nil {
		// Добавляем контекст к ошибке
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

// Begin - мок для метода Begin.
func (m *MockDatabase) Begin(ctx context.Context) (pgx.Tx, error) {
	a := m.Called(ctx)

	// Проверка типа
	tx, ok := a.Get(0).(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("failed to assert type for pgx.Tx: got %T", a.Get(0))
	}

	// Обработка и добавление контекста к ошибке
	if err := a.Error(1); err != nil {
		return tx, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return tx, nil
}

// Close - мок для метода Close.
func (m *MockDatabase) Close() {
	m.Called()
}
