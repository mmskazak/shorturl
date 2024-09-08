package postgresql

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/storage/postgresql/mocks"
	"testing"
)

func TestPostgreSQL_DeleteURLs(t *testing.T) {
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
