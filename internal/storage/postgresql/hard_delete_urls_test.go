package postgresql

import (
	"context"
	"errors"
	"testing"
	"time"

	"mmskazak/shorturl/internal/storage/postgresql/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func Test_hardDeleteSoftDeletedURLs_Success(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	pool := new(mocks.MockDatabase)
	zapLog := zap.NewNop().Sugar()

	timeSleep = time.Second * 1
	pool.On("Exec", ctx, mock.Anything, mock.Anything).
		Return(pgconn.NewCommandTag("DELETE 0"), nil)
	hardDeleteSoftDeletedURLs(ctx, pool, zapLog)
	pool.AssertExpectations(t)
}

func Test_hardDeleteSoftDeletedURLs_Err(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	pool := new(mocks.MockDatabase)
	zapLog := zap.NewNop().Sugar()

	timeSleep = time.Second * 1
	pool.On("Exec", ctx, mock.Anything, mock.Anything).
		Return(pgconn.NewCommandTag(""), errors.New("test error"))
	hardDeleteSoftDeletedURLs(ctx, pool, zapLog)
	pool.AssertExpectations(t)
}
