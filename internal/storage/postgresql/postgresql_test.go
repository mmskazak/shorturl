package postgresql

import (
	"context"
	"errors"
	"testing"

	"mmskazak/shorturl/internal/storage/postgresql/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPostgreSQL_Ping_NoError(t *testing.T) {
	ctx := context.Background()
	mockPool := new(mocks.MockDatabase)

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}
	mockPool.On("Ping", ctx).Return(nil)
	err := s.Ping(ctx)
	assert.NoError(t, err)
}

func TestPostgreSQL_Ping_ErrPing(t *testing.T) {
	ctx := context.Background()
	mockPool := new(mocks.MockDatabase)

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}
	mockPool.On("Ping", ctx).Return(errors.New("test error"))
	err := s.Ping(ctx)
	assert.Error(t, err)
}

func TestPostgreSQL_Close_Success(t *testing.T) {
	mockPool := new(mocks.MockDatabase)

	// Создаем объект PostgreSQL с моками
	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}
	mockPool.On("Close").Once()
	err := s.Close()
	require.NoError(t, err)
	mockPool.AssertExpectations(t)
}

func TestPostgreSQL_Close_ErrPoolNil(t *testing.T) {
	dsn := ""
	zapLog := zap.NewNop().Sugar()
	err := runMigrations(dsn, zapLog)
	assert.EqualError(t, err, "error opening migrations directory: URL cannot be empty")
}
