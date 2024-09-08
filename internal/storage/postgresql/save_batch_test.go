package postgresql

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/models"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"
)

func TestPostgreSQL_SaveBatch_BeginTxError(t *testing.T) {
	ctx := context.Background()
	incoming := []models.Incoming{
		{
			CorrelationID: "123",
			OriginalURL:   "http://yandex.ru",
		},
		{
			CorrelationID: "456",
			OriginalURL:   "http://google.com",
		},
	}
	baseHost := "http://localhost"
	userID := "1"
	generator := genidurl.NewGenIDService()

	mockPool := new(mocks.MockDatabase)
	mockTx := new(mocks.MockTx)

	s := &PostgreSQL{
		pool:   mockPool,
		zapLog: zap.NewNop().Sugar(),
	}

	mockPool.On("Begin", ctx).Return(mockTx, errors.New("test error")).Once()

	got, err := s.SaveBatch(ctx, incoming, baseHost, userID, generator)
	assert.Error(t, err)
	assert.Equal(t, []models.Output(nil), got)

	mockPool.AssertExpectations(t)
}
