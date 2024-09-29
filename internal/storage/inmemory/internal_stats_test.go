package inmemory

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"mmskazak/shorturl/internal/models"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInMemory_InternalStats(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		mu        *sync.Mutex
		data      map[string]models.URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	tests := []struct {
		name    string
		fields  fields
		want    models.Stats
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty storage",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]models.URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			want:    models.Stats{Urls: strconv.Itoa(0), Users: strconv.Itoa(0)},
			wantErr: assert.NoError, // Ожидаем, что ошибки не будет
		},
		{
			name: "with records and users",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]models.URLRecord{
					"1": {},
					"2": {},
					"3": {},
				},
				userIndex: map[string][]string{
					"user1": {"1"},
					"user2": {"2", "3"},
				},
				zapLog: zap.NewNop().Sugar(),
			},
			want:    models.Stats{Urls: strconv.Itoa(3), Users: strconv.Itoa(2)},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InMemory{
				mu:        tt.fields.mu,
				data:      tt.fields.data,
				userIndex: tt.fields.userIndex,
				zapLog:    tt.fields.zapLog,
			}
			got, err := m.InternalStats(ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("InternalStats(%v)", ctx)) {
				return
			}
			assert.Equalf(t, tt.want, got, "InternalStats(%v)", ctx)
		})
	}
}
