package postgresql

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
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
