package api

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"mmskazak/shorturl/internal/contracts/mocks"
	"mmskazak/shorturl/internal/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestInternalStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	zapLog := zaptest.NewLogger(t).Sugar()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", http.NoBody)

	expectedResult := models.Stats{
		Urls:  strconv.Itoa(300),
		Users: strconv.Itoa(10),
	}
	store := mocks.NewMockIInternalStats(ctrl)
	store.EXPECT().
		InternalStats(ctx).
		Return(expectedResult, nil) // Замените expectedResult на то, что должен возвращать метод

	InternalStats(ctx, w, req, store, zapLog)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"urls\":\"300\",\"users\":\"10\"}", w.Body.String())
}
