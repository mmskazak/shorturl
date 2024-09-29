package inmemory

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/models"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInMemory_InternalStats(t *testing.T) {
	type fields struct {
		mu        *sync.Mutex
		data      map[string]models.URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			args: args{
				in0: context.Background(),
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
			args: args{
				in0: context.Background(),
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
			got, err := m.InternalStats(tt.args.in0)
			if !tt.wantErr(t, err, fmt.Sprintf("InternalStats(%v)", tt.args.in0)) {
				return
			}
			assert.Equalf(t, tt.want, got, "InternalStats(%v)", tt.args.in0)
		})
	}
}
