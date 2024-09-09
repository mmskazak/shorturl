package postgresql

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"
)

func TestPostgreSQL_DeleteURLs_EmptyURLs(t *testing.T) {
	mockPool := new(mocks.MockDatabase)
	ctx := context.Background()
	var urlIDs []string
	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}
	err := s.DeleteURLs(ctx, urlIDs)
	assert.Equal(t, nil, err)
}

func TestPostgreSQL_DeleteURLs_ErrBegin(t *testing.T) {
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	ctx := context.Background()

	urlIDs := []string{
		"testtest",
		"test1234",
		"1234test",
	}

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}
	mockPool.On("Begin", ctx).Return(mockTx, errors.New("test error"))
	err := s.DeleteURLs(ctx, urlIDs)
	assert.Equal(t, "failed to begin transaction: test error", err.Error())
}

func TestPostgreSQL_DeleteURLs_Success(t *testing.T) {
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	mockBatchResults := new(mocks.MockBatchResults)
	ctx := context.Background()

	urlIDs := []string{
		"testtest",
		"test1234",
		"1234test",
	}

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}
	mockPool.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Rollback", ctx).Return(nil)
	mockTx.On("SendBatch", ctx, mock.Anything).Return(mockBatchResults)
	mockBatchResults.On("Close").Return(nil)
	mockTx.On("Commit", ctx).Return(nil)
	err := s.DeleteURLs(ctx, urlIDs)
	assert.NoError(t, err)
}

func TestPostgreSQL_DeleteURLs_CommitError(t *testing.T) {
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	mockBatchResults := new(mocks.MockBatchResults)
	ctx := context.Background()

	urlIDs := []string{
		"testtest",
		"test1234",
		"1234test",
	}

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}
	mockPool.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Rollback", ctx).Return(nil)
	mockTx.On("SendBatch", ctx, mock.Anything).Return(mockBatchResults)
	mockBatchResults.On("Close").Return(nil)
	mockTx.On("Commit", ctx).Return(errors.New("test error"))
	err := s.DeleteURLs(ctx, urlIDs)
	require.Error(t, err)
	assert.EqualError(t, err, "failed to commit transaction: test error")
}

func TestPostgreSQL_DeleteURLs_CloseError(t *testing.T) {
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	mockBatchResults := new(mocks.MockBatchResults)
	ctx := context.Background()

	urlIDs := []string{
		"testtest",
		"test1234",
		"1234test",
	}

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}
	mockPool.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Rollback", ctx).Return(nil)
	mockTx.On("SendBatch", ctx, mock.Anything).Return(mockBatchResults)
	mockBatchResults.On("Close").Return(errors.New("test error"))

	err := s.DeleteURLs(ctx, urlIDs)
	require.Error(t, err)
	assert.EqualError(t, err, "failed to delete URLs in batch: test error")
}
