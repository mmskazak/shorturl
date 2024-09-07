package inmemory

import (
	"context"
	"sync"
	"testing"

	"mmskazak/shorturl/internal/models"

	"go.uber.org/zap"
)

func TestInMemory_SetShortURL(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		mu        *sync.Mutex
		data      map[string]models.URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	type args struct {
		id          string
		originalURL string
		userID      string
		deleted     bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty storage, add new URL",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]models.URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{

				id:          "short1",
				originalURL: "http://example.com",
				userID:      "user1",
				deleted:     false,
			},
			wantErr: false,
		},
		{
			name: "err key already exists",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]models.URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://old.com", UserID: "user1", Deleted: false},
				},
				userIndex: map[string][]string{
					"user1": {"short1"},
				},
				zapLog: zap.NewNop().Sugar(),
			},
			args: args{
				id:          "short1",
				originalURL: "http://new.com",
				userID:      "user1",
				deleted:     false,
			},
			wantErr: true,
		},
		{
			name: "add URL for a new user",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]models.URLRecord{
					"short1": {ShortURL: "short1", OriginalURL: "http://example.com", UserID: "user1", Deleted: false},
				},
				userIndex: map[string][]string{
					"user1": {"short1"},
				},
				zapLog: zap.NewNop().Sugar(),
			},
			args: args{
				id:          "short2",
				originalURL: "http://example2.com",
				userID:      "user2",
				deleted:     false,
			},
			wantErr: false,
		},
		{
			name: "add deleted URL",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]models.URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				id:          "short1",
				originalURL: "http://example.com",
				userID:      "user1",
				deleted:     true,
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
			err := m.SetShortURL(ctx, tt.args.id, tt.args.originalURL, tt.args.userID, tt.args.deleted)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetShortURL() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				// Проверка, что URL был добавлен в data
				if record, exists := m.data[tt.args.id]; !exists || record.OriginalURL != tt.args.originalURL ||
					record.UserID != tt.args.userID || record.Deleted != tt.args.deleted {
					t.Errorf("SetShortURL() did not store the correct URLRecord, got = %v, want = %v",
						record, tt.args)
				}

				// Проверка, что URL был добавлен в индекс пользователя
				if _, exists := m.userIndex[tt.args.userID]; !exists {
					t.Errorf("SetShortURL() did not add URL to userIndex for user %v", tt.args.userID)
				}
			}
		})
	}
}
