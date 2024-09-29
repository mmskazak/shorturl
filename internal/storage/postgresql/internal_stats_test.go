package postgresql

import (
	"context"
	"mmskazak/shorturl/internal/models"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPostgreSQL_InternalStats(t *testing.T) {
	// Настройка моков
	mockDB := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)
	mockRowUrls := new(mocks.MockRow)
	mockRowUsers := new(mocks.MockRow)
	ctx := context.Background()
	zapLog := zap.NewNop().Sugar()

	// Инициализация PostgreSQL с использованием мока
	pg := &PostgreSQL{
		pool:   mockDB,
		zapLog: zapLog,
	}

	// Определение ожидаемого поведения
	mockDB.On("Begin", ctx).Return(mockTx, nil)

	mockTx.On("QueryRow", ctx, "SELECT COUNT(original_url) FROM urls WHERE 1").
		Return(mockRowUrls)
	mockRowUrls.On("Scan", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		urls, ok := args.Get(0).(*string)
		if ok {
			*urls = "5"
		}
	}).Return(nil)

	mockTx.On("QueryRow", ctx, "SELECT COUNT(DISTINCT user_id) FROM urls").
		Return(mockRowUsers)
	mockRowUsers.On("Scan", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		users, ok := args.Get(1).(*string)
		if ok {
			*users = "2"
		}
	}).Return(nil)

	mockTx.On("Rollback", ctx).Return(nil)

	// Вызов функции
	stats, err := pg.InternalStats(ctx)

	// Проверка результатов
	assert.NoError(t, err)
	assert.Equal(t, models.Stats{Urls: "5", Users: "2"}, stats)

	// Проверка, что все ожидания были выполнены
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockRowUrls.AssertExpectations(t)
	mockRowUsers.AssertExpectations(t)
}
