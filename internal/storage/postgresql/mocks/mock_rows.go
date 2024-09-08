package mocks

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockRows представляет мок для интерфейса Rows
type MockRows struct {
	mock.Mock
}

// Close Реализация метода
func (m *MockRows) Close() {
	m.Called()
}

// Err Реализация метода
func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

// CommandTag Реализация метода
func (m *MockRows) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

// FieldDescriptions Реализация метода
func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	return args.Get(0).([]pgconn.FieldDescription)
}

// Next Реализация метода
func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

// Scan Реализация метода
func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

// Values Реализация метода
func (m *MockRows) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

// RawValues Реализация метода
func (m *MockRows) RawValues() [][]byte {
	args := m.Called()
	return args.Get(0).([][]byte)
}

// Conn Реализация метода
func (m *MockRows) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}
