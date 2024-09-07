package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"
)

func TestPostgreSQL_SetShortURL_ErrBeginTx(t *testing.T) {
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	ctx := context.Background()
	shortURL := "testtest"
	targetURL := "http://google.com"
	userID := "1"
	deleted := false

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}
	mockPool.On("Begin", ctx).Return(mockTx, errors.New("test error")).Once()

	err := s.SetShortURL(ctx, shortURL, targetURL, userID, deleted)
	assert.EqualError(t, err, "error beginning transaction: test error")

	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostgreSQL_SetShortURL_Success(t *testing.T) {
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	ctx := context.Background()
	shortURL := "testtest"
	targetURL := "http://google.com"
	userID := "1"
	deleted := false

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}

	commandTag := pgconn.NewCommandTag("INSERT 0 1")

	mockPool.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Exec", ctx, mock.AnythingOfType("string"),
		[]interface{}{shortURL, targetURL, userID, deleted}).
		Return(commandTag, nil).Once()
	mockTx.On("Commit", ctx).Return(nil)
	mockTx.On("Rollback", ctx).Return(nil) // Добавлено ожидание для Rollback

	err := s.SetShortURL(ctx, shortURL, targetURL, userID, deleted)
	assert.NoError(t, err)

	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
