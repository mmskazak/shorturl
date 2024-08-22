package inmemory

import (
	"context"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestInMemory_GetShortURL(t *testing.T) {
	type fields struct {
		mu        *sync.Mutex
		data      map[string]URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	type args struct {
		in0 context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "successful retrieval",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com", UserID: "user1", Deleted: false},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				in0: context.Background(),
				id:  "short1",
			},
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name: "URL not found",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				in0: context.Background(),
				id:  "short1",
			},
			want:    "",
			wantErr: true, // Ожидаем ошибку, так как URL не найден
		},
		{
			name: "URL is deleted",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com", UserID: "user1", Deleted: true},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				in0: context.Background(),
				id:  "short1",
			},
			want:    "",
			wantErr: true, // Ожидаем ошибку или пустое значение, так как URL помечен как удаленный
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
			got, err := m.GetShortURL(tt.args.in0, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetShortURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
