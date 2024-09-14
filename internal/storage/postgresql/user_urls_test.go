package postgresql

import (
	"context"
	"errors"
	"testing"

	"mmskazak/shorturl/internal/models"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPostgreSQL_GetUserURLs_ErrorQuery(t *testing.T) {
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

func TestPostgreSQL_GetUserURLs_ErrScan(t *testing.T) {
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
		Return(mockRows, nil).Once()
	mockRows.On("Next").Return(true).Once() // Первый вызов вернет true
	mockRows.On("Scan", mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
		Return(errors.New("test error"))
	mockRows.On("Close").Return(nil).Once()

	var expected []models.URL

	// Вызываем тестируемую функцию
	got, err := s.GetUserURLs(ctx, userID, baseHost)
	assert.EqualError(t, err, "error scanning row: test error")
	assert.Equal(t, expected, got)

	// Проверяем вызовы методов моков
	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestPostgreSQL_GetUserURLs_BeginErr(t *testing.T) {
	// Создаем моки
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)

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
	mockPool.On("Begin", ctx).Return(mockTx, errors.New("test error")).Once()

	// Вызываем тестируемую функцию
	got, err := s.GetUserURLs(ctx, userID, baseHost)
	assert.Error(t, err)
	assert.Equal(t, []models.URL(nil), got)

	// Проверяем вызовы методов моков
	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostgreSQL_GetUserURLs(t *testing.T) {
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
		Return(mockRows, nil).Once()
	mockRows.On("Next").Return(true).Once()  // Первый вызов вернет true
	mockRows.On("Next").Return(false).Once() // Второй вызов вернет false
	mockRows.On("Err").Return(nil).Once()
	mockRows.On("Scan", mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
		Run(func(args mock.Arguments) {
			shortURL, ok := args.Get(0).(*string)
			if ok {
				*shortURL = "testtest"
			}
			originalURL, ok := args.Get(1).(*string)
			if ok {
				*originalURL = "http://google.com"
			}
		}).Return(nil)
	mockTx.On("Commit", ctx).Return(nil).Once()
	mockRows.On("Close").Return(nil).Once()

	expected := []models.URL{
		{
			ShortURL:    "http://localhost/testtest",
			OriginalURL: "http://google.com",
		},
	}

	// Вызываем тестируемую функцию
	got, err := s.GetUserURLs(ctx, userID, baseHost)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	// Проверяем вызовы методов моков
	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}
