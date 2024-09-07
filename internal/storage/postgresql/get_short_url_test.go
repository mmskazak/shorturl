package postgresql

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"
)

func TestPostgreSQL_GetShortURL(t *testing.T) {
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
