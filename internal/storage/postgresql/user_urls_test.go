package postgresql

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"
)

func TestPostgreSQL_GetUserURLs(t *testing.T) {
	t.Skip()
	// Создаем моки
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	mockRows := new(mocks.MockRows)

	// Определяем контекст и входные данные
	ctx := context.Background()
	userID := "1"
	baseHost := "http://localhost"

	// Создаем объект PostgreSQL с моками
	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}

	// Устанавливаем ожидаемые вызовы методов моков
	mockPool.On("Begin", ctx).Return(mockTx, nil).Once()
	mockTx.On("Rollback", ctx).Return(nil)
	mockTx.On("Query", ctx, mock.Anything, []interface{}{userID}).Return(mockRows, nil).Once()
	mockRows.On("Next").Return(true).Once() // Будем возвращать одну строку
	mockRows.On("Scan", mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
		Run(func(args mock.Arguments) {
			shortURLPtr := args.Get(0).(*string)
			*shortURLPtr = "/short"
			originalURLPtr := args.Get(1).(*string)
			*originalURLPtr = "http://example.com/original"
		}).Return(nil)
	mockRows.On("Next").Return(false)
	mockRows.On("Close").Return(nil)
	mockRows.On("Err").Return(nil)
	mockTx.On("Commit", ctx).Return(nil)

	// Вызываем тестируемую функцию
	got, err := s.GetUserURLs(ctx, userID, baseHost)

	// Проверяем ошибки
	require.NoError(t, err)
	require.Len(t, got, 1, "Expected 1 URL, got %d", len(got))

	// Проверяем корректность полученного URL
	expectedShortURL := "http://localhost/short"
	assert.Equal(t, expectedShortURL, got[0].ShortURL)
	assert.Equal(t, "http://example.com/original", got[0].OriginalURL)

	// Проверяем вызовы методов моков
	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}
