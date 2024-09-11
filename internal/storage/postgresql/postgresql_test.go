package postgresql

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"
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
	// Создаем объект PostgreSQL с моками
	s := &PostgreSQL{
		pool:   nil,
		zapLog: zap.NewNop().Sugar(),
	}

	err := s.Close()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pool is nil")
}
