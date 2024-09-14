package postgresql

import (
	"context"
	"errors"
	"testing"

	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPostgreSQL_GetShortURL_Success(t *testing.T) {
	pool := new(mocks.MockDatabase)
	mockRow := new(mocks.MockRow)

	s := &PostgreSQL{
		pool:   pool,
		zapLog: zap.NewNop().Sugar(),
	}

	shortURL := "testtest"

	pool.On("QueryRow",
		context.Background(),
		"SELECT original_url, deleted FROM urls WHERE short_url = $1",
		[]interface{}{shortURL}, // Используем слайс вместо строки
	).Return(mockRow, nil)
	mockRow.On("Scan", mock.AnythingOfType("*string"), mock.AnythingOfType("*bool")).
		Run(func(args mock.Arguments) {
			// Проверяем приведение типов для первого аргумента
			strPtr, ok := args.Get(0).(*string)
			if ok {
				*strPtr = "http://google.com" // Присваиваем значение, если приведение типа прошло успешно
			}

			// Проверяем приведение типов для второго аргумента
			boolPtr, ok := args.Get(1).(*bool)
			if ok {
				*boolPtr = false // Присваиваем значение, если приведение типа прошло успешно
			}
		}).
		Return(nil)

	got, err := s.GetShortURL(context.Background(), "testtest")
	require.NoError(t, err)
	assert.Equal(t, "http://google.com", got)
}

func TestPostgreSQL_GetShortURL_ErrNoRows(t *testing.T) {
	pool := new(mocks.MockDatabase)
	mockRow := new(mocks.MockRow)

	s := &PostgreSQL{
		pool:   pool,
		zapLog: zap.NewNop().Sugar(),
	}

	shortURL := "testtest"

	pool.On("QueryRow",
		context.Background(),
		"SELECT original_url, deleted FROM urls WHERE short_url = $1",
		[]interface{}{shortURL}, // Используем слайс вместо строки
	).Return(mockRow, nil)
	mockRow.On("Scan", mock.AnythingOfType("*string"), mock.AnythingOfType("*bool")).
		Return(pgx.ErrNoRows)

	got, err := s.GetShortURL(context.Background(), "testtest")
	assert.Error(t, pgx.ErrNoRows, err)
	assert.Equal(t, "", got)
}

func TestPostgreSQL_GetShortURL_SomeErr(t *testing.T) {
	pool := new(mocks.MockDatabase)
	mockRow := new(mocks.MockRow)

	s := &PostgreSQL{
		pool:   pool,
		zapLog: zap.NewNop().Sugar(),
	}

	shortURL := "testtest"
	testError := errors.New("test error")
	pool.On("QueryRow",
		context.Background(),
		"SELECT original_url, deleted FROM urls WHERE short_url = $1",
		[]interface{}{shortURL}, // Используем слайс вместо строки
	).Return(mockRow, nil)
	mockRow.On("Scan", mock.AnythingOfType("*string"), mock.AnythingOfType("*bool")).
		Return(errors.New("test error"))

	got, err := s.GetShortURL(context.Background(), "testtest")
	assert.Error(t, testError, err)
	assert.Equal(t, "", got)
}

func TestPostgreSQL_GetShortURL_Deleted(t *testing.T) {
	pool := new(mocks.MockDatabase)
	mockRow := new(mocks.MockRow)

	s := &PostgreSQL{
		pool:   pool,
		zapLog: zap.NewNop().Sugar(),
	}

	shortURL := "testtest"

	pool.On("QueryRow",
		context.Background(),
		"SELECT original_url, deleted FROM urls WHERE short_url = $1",
		[]interface{}{shortURL}, // Используем слайс вместо строки
	).Return(mockRow, nil)
	mockRow.On("Scan", mock.AnythingOfType("*string"), mock.AnythingOfType("*bool")).
		Run(func(args mock.Arguments) {
			// Проверяем приведение типов для первого аргумента
			strPtr, ok := args.Get(0).(*string)
			if ok {
				*strPtr = "http://google.com" // Присваиваем значение, если приведение типа прошло успешно
			}

			// Проверяем приведение типов для второго аргумента
			boolPtr, ok := args.Get(1).(*bool)
			if ok {
				*boolPtr = true // Присваиваем значение, если приведение типа прошло успешно
			}
		}).
		Return(nil)

	got, err := s.GetShortURL(context.Background(), "testtest")
	assert.Error(t, storageErrors.ErrDeletedShortURL, err)
	assert.Equal(t, "", got)
}
