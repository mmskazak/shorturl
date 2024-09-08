package postgresql

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/models"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"
)

func TestPostgreSQL_GetUserURLs(t *testing.T) {
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
	mockTx.On("Query", ctx, mock.Anything, []interface{}{userID}).
		Return(mockRows, errors.New("test error")).Once()

	// Вызываем тестируемую функцию
	got, err := s.GetUserURLs(ctx, userID, baseHost)
	assert.Error(t, err)
	assert.Equal(t, []models.URL(nil), got)

	// Проверяем вызовы методов моков
	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}
