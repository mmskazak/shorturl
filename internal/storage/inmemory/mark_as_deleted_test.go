package inmemory

import (
	"context"
	"mmskazak/shorturl/internal/models"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestInMemory_MarkURLAsDeleted(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		mu        *sync.Mutex
		data      map[string]models.URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful deletion",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]models.URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com", UserID: "user1", Deleted: false},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				id: "short1",
			},
			wantErr: false,
		},
		{
			name: "URL not found",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord), // Пустое хранилище
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				id: "short2",
			},
			wantErr: true, // Ожидаем ошибку, так как URL нет в хранилище
		},
		{
			name: "repeated deletion",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]models.URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com", UserID: "user1", Deleted: true},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				id: "short1",
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
			err := m.MarkURLAsDeleted(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkURLAsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Дополнительная проверка для успешного удаления
			if !tt.wantErr {
				record, exists := m.data[tt.args.id]
				if !exists {
					t.Errorf("MarkURLAsDeleted() did not find the record after deletion, id = %v", tt.args.id)
				}
				if !record.Deleted {
					t.Errorf("MarkURLAsDeleted() did not mark the URL as deleted, id = %v", tt.args.id)
				}
			}
		})
	}
}
