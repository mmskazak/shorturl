package postgresql

import (
	"context"
	"errors"
	"testing"

	"mmskazak/shorturl/internal/models"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPostgreSQL_SaveBatch_BeginTxError(t *testing.T) {
	ctx := context.Background()
	incoming := []models.Incoming{
		{
			CorrelationID: "123",
			OriginalURL:   "http://yandex.ru",
		},
		{
			CorrelationID: "456",
			OriginalURL:   "http://google.com",
		},
	}
	baseHost := "http://localhost"
	userID := "1"
	generator := genidurl.NewGenIDService()

	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}

	mockPool.On("Begin", ctx).Return(mockTx, errors.New("test error")).Once()

	got, err := s.SaveBatch(ctx, incoming, baseHost, userID, generator)
	assert.Error(t, err)
	assert.Equal(t, []models.Output(nil), got)

	mockPool.AssertExpectations(t)
}

func TestPostgreSQL_SaveBatch_ErrIncoming(t *testing.T) {
	ctx := context.Background()
	var incoming []models.Incoming
	baseHost := "http://localhost"
	userID := "1"
	generator := genidurl.NewGenIDService()

	mockPool := new(mocks.MockDatabase)

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}

	got, err := s.SaveBatch(ctx, incoming, baseHost, userID, generator)
	assert.Error(t, err)
	assert.Equal(t, []models.Output(nil), got)
}

func TestPostgreSQL_SaveBatch(t *testing.T) {
	ctx := context.Background()
	incoming := []models.Incoming{
		{
			CorrelationID: "123",
			OriginalURL:   "http://yandex.ru",
		},
		{
			CorrelationID: "456",
			OriginalURL:   "http://google.com",
		},
	}
	baseHost := "http://localhost"
	userID := "1"
	generator := genidurl.NewGenIDService()

	mockTx := new(mocks.MockTx)
	mockPool := new(mocks.MockDatabase)
	mockBatchResults := new(mocks.MockBatchResults)
	mockRow := new(mocks.MockRow)
	mockPool.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("SendBatch", ctx, mock.Anything).Return(mockBatchResults)
	mockTx.On("Commit", ctx).Return(nil)
	mockTx.On("Rollback", ctx).Return(nil)

	callCount := 0
	// Настройка мока для метода Scan
	mockRow.On("Scan", mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
		Run(func(args mock.Arguments) {
			callCount++
			shortURL, ok := args.Get(0).(*string)
			if ok {
				if callCount == 1 {
					*shortURL = "testt123" // Значение для первого вызова
				} else if callCount == 2 {
					*shortURL = "testt456" // Значение для второго вызова
				}
			}

			originalURL, ok := args.Get(1).(*string)
			if ok {
				if callCount == 1 {
					*originalURL = "http://yandex.ru"
				} else if callCount == 2 {
					*originalURL = "http://google.com"
				}
			}
		}).Return(nil).Times(2) // Указываем, что ожидаем два вызова

	mockBatchResults.On("QueryRow").Return(mockRow)
	mockBatchResults.On("Exec").Return(pgconn.NewCommandTag("INSERT 0 1"), nil)
	mockBatchResults.On("Close").Return(nil)

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}

	output := []models.Output{
		models.Output{CorrelationID: "123", ShortURL: "http://localhost/testt123"},
		models.Output{CorrelationID: "456", ShortURL: "http://localhost/testt456"},
	}

	got, err := s.SaveBatch(ctx, incoming, baseHost, userID, generator)
	assert.Equal(t, output, got)
	assert.NoError(t, err)
}
