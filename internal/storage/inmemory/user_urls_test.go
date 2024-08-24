package inmemory

import (
	"context"
	"mmskazak/shorturl/internal/models"
	"reflect"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestInMemory_GetUserURLs(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		mu        *sync.Mutex
		data      map[string]URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	type args struct {
		userID   string
		baseHost string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.URL
		wantErr bool
	}{
		{
			name: "empty data",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
			},
			args: args{

				userID:   "11111",
				baseHost: "http://localhost",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "user with URLs",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com/1", UserID: "11111", Deleted: false},
					"short2": {ShortURL: "short2", OriginalURL: "http://example.com/2", UserID: "11111", Deleted: false},
				},
				userIndex: map[string][]string{
					"11111": {"short1", "short2"},
				},
				zapLog: zap.NewNop().Sugar(),
			},
			args: args{
				userID:   "11111",
				baseHost: "http://localhost",
			},
			want: []models.URL{
				{OriginalURL: "http://example.com/1", ShortURL: "http://localhost/short1"},
				{OriginalURL: "http://example.com/2", ShortURL: "http://localhost/short2"},
			},
			wantErr: false,
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
			got, err := m.GetUserURLs(ctx, tt.args.userID, tt.args.baseHost)
			if got == nil {
				got = []models.URL{}
			}
			if tt.want == nil {
				tt.want = []models.URL{}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserURLs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
