package postgresql

import (
	"context"
	"errors"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

func TestPostgreSQL_SetShortURL_SomeErrExec(t *testing.T) {
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

	commandTag := pgconn.NewCommandTag("INSERT 0 0")

	mockPool.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Exec", ctx, mock.AnythingOfType("string"),
		[]interface{}{shortURL, targetURL, userID, deleted}).
		Return(commandTag, errors.New("test error")).Once()
	mockTx.On("Rollback", ctx).Return(nil) // Добавлено ожидание для Rollback

	err := s.SetShortURL(ctx, shortURL, targetURL, userID, deleted)
	assert.Error(t, err)

	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostgreSQL_SetShortURL_ErrUniqueOriginalUrl(t *testing.T) {
	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	mockRow := new(mocks.MockRow)
	ctx := context.Background()
	shortURL := "testtest"
	targetURL := "http://google.com"
	userID := "1"
	deleted := false

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}

	commandTag := pgconn.NewCommandTag("INSERT 0 0")

	mockPool.On("Begin", ctx).Return(mockTx, nil)

	pgErr := &pgconn.PgError{
		Code:           pgerrcode.UniqueViolation,
		ConstraintName: "unique_original_url",
	}

	mockTx.On("Exec", ctx, mock.AnythingOfType("string"),
		[]interface{}{shortURL, targetURL, userID, deleted}).
		Return(commandTag, pgErr)

	mockPool.On("QueryRow", ctx, "SELECT short_url FROM urls WHERE original_url = $1",
		[]interface{}{"http://google.com"}).Return(mockRow)

	mockRow.On("Scan", mock.AnythingOfType("*string")).
		Run(func(args mock.Arguments) {
			strPtr, ok := args.Get(0).(*string)
			if ok {
				*strPtr = "gogotest"
			}
		}).
		Return(nil)

	mockTx.On("Rollback", ctx).Return(nil) // Добавлено ожидание для Rollback

	err := s.SetShortURL(ctx, shortURL, targetURL, userID, deleted)
	require.Error(t, err)
	assert.EqualError(t, err, "error original url already exists: gogotest")

	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestPostgreSQL_SetShortURL_ErrUniqueShortUrl(t *testing.T) {
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

	commandTag := pgconn.NewCommandTag("INSERT 0 0")

	mockPool.On("Begin", ctx).Return(mockTx, nil)

	pgErr := &pgconn.PgError{
		Code:           pgerrcode.UniqueViolation,
		ConstraintName: "unique_short_url",
	}

	mockTx.On("Exec", ctx, mock.AnythingOfType("string"),
		[]interface{}{shortURL, targetURL, userID, deleted}).
		Return(commandTag, pgErr)

	mockTx.On("Rollback", ctx).Return(nil) // Добавлено ожидание для Rollback

	err := s.SetShortURL(ctx, shortURL, targetURL, userID, deleted)
	require.Error(t, err)
	assert.EqualError(t, err, "error key already exists")

	mockPool.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
