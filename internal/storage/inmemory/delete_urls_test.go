package inmemory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"mmskazak/shorturl/internal/models"

	"go.uber.org/zap"
)

func TestInMemory_DeleteURLs(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		mu        *sync.Mutex
		data      map[string]models.URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	type args struct {
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
				data: map[string]models.URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com", UserID: "user1", Deleted: false},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				urlIDs: []string{"short1", "short2"},
			},
			wantErr: false, // Ожидаем отсутствие ошибки, даже если один из URL не существует
		},
		{
			name: "empty list of URLs",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]models.URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				urlIDs: []string{},
			},
			wantErr: false, // Ожидаем отсутствие ошибки при пустом списке
		},
		{
			name: "all URLs are already deleted",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]models.URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com", UserID: "user1", Deleted: true},
					"short2": {ShortURL: "short2", OriginalURL: "http://example2.com", UserID: "user1", Deleted: true},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
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
			err := m.DeleteURLs(ctx, tt.args.urlIDs)
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

func Test_splitSlice_chunkSize(t *testing.T) {
	input := []string{
		"one",
		"two",
	}
	chunkSize := 0
	got := splitSlice(input, chunkSize)
	assert.Nil(t, got)
}

func Test_splitSlice_input(t *testing.T) {
	var input []string
	chunkSize := 5000
	got := splitSlice(input, chunkSize)
	assert.Equal(t, [][]string{}, got)
}
