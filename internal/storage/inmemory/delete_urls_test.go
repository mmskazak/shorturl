package inmemory

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"testing"
	"time"
)

func TestInMemory_DeleteURLs(t *testing.T) {
	type fields struct {
		mu        *sync.Mutex
		data      map[string]URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	type args struct {
		ctx    context.Context
		urlIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "some URLs do not exist",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com", UserID: "user1", Deleted: false},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				ctx:    context.Background(),
				urlIDs: []string{"short1", "short2"},
			},
			wantErr: false, // Ожидаем отсутствие ошибки, даже если один из URL не существует
		},
		{
			name: "empty list of URLs",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				ctx:    context.Background(),
				urlIDs: []string{},
			},
			wantErr: false, // Ожидаем отсутствие ошибки при пустом списке
		},
		{
			name: "all URLs are already deleted",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com", UserID: "user1", Deleted: true},
					"short2": {ShortURL: "short2", OriginalURL: "http://example2.com", UserID: "user1", Deleted: true},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				ctx:    context.Background(),
				urlIDs: []string{"short1", "short2"},
			},
			wantErr: false, // Ожидаем отсутствие ошибки, так как URL уже удалены
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
			err := m.DeleteURLs(tt.args.ctx, tt.args.urlIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteURLs() error = %v, wantErr %v", err, tt.wantErr)
			}
			time.Sleep(1 * time.Second) // Ждем секунду, горутину которая помечает как удаленные

			// Дополнительная проверка для успешного удаления
			if !tt.wantErr {
				for _, id := range tt.args.urlIDs {
					record, exists := m.data[id]
					if exists && !record.Deleted {
						t.Errorf("DeleteURLs() did not mark URL as deleted, id = %v", id)
					}
				}
			}
		})
	}
}
