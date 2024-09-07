package mocks

import "github.com/stretchr/testify/mock"

// MockRow — структура, представляющая мок для интерфейса Row
type MockRow struct {
	mock.Mock
}

// Реализация метода Scan для мока.
func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	// Присвоение значений, если они были установлены через вызов Run
	if len(args) > 1 && args.Get(0) == nil {
		*(dest[0].(*string)) = args.String(1)
		*(dest[1].(*bool)) = args.Bool(2)
	}
	return args.Error(0)
}
