package postgresql

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"mmskazak/shorturl/internal/storage/postgresql/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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
	assert.Error(t, err)
	assert.Equal(t, "failed to begin transaction: failed to begin transaction: test error", err.Error())
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
	assert.EqualError(t, err, "failed to commit transaction: error MockTx func Commit: test error")
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
	assert.EqualError(t, err, "failed to delete URLs in batch: failed to close mock batch results: test error")
}

func TestPostgreSQL_DeleteURLs_batchSize5000(t *testing.T) {
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	mockBatchResults := new(mocks.MockBatchResults)
	ctx := context.Background()

	urlIDs := make([]string, 0, 5000)
	for range 5000 {
		urlIDs = append(urlIDs, "TestTest")
	}
	fmt.Println(len(urlIDs))
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
	assert.EqualError(t, err, "failed to delete URLs in batch: failed to close mock batch results: test error")
}
