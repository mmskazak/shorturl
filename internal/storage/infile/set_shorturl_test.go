package infile

import (
	"context"
	"testing"

	"mmskazak/shorturl/internal/storage/inmemory"

	"go.uber.org/zap"
)

func TestInFile_SetShortURL(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		InMe     *inmemory.InMemory
		zapLog   *zap.SugaredLogger
		filePath string
	}
	type args struct {
		idShortPath string
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
			name: "Successful URL addition",
			fields: fields{
				InMe: func() *inmemory.InMemory {
					in, _ := inmemory.NewInMemory(zap.NewNop().Sugar())
					return in
				}(), // инициализируем InMemory
				zapLog:   zap.NewNop().Sugar(),
				filePath: "/path/to/storage",
			},
			args: args{

				idShortPath: "short123",
				originalURL: "https://example.com",
				userID:      "user1",
				deleted:     false,
			},
			wantErr: false,
		},
		{
			name: "Duplicate short URL",
			fields: fields{
				InMe: func() *inmemory.InMemory {
					in, _ := inmemory.NewInMemory(zap.NewNop().Sugar())
					err := in.SetShortURL(context.Background(),
						"short123",
						"https://example.com",
						"1",
						false)
					if err != nil {
						return nil
					}
					return in
				}(), // инициализируем InMemory
				zapLog:   zap.NewNop().Sugar(),
				filePath: "/path/to/storage",
			},
			args: args{
				idShortPath: "short123", // Дублируем тот же short URL
				originalURL: "https://test.com",
				userID:      "user1",
				deleted:     false,
			},
			wantErr: true, // Ожидаем ошибку из-за дублирования
		},
		{
			name: "Empty original URL",
			fields: fields{
				InMe: func() *inmemory.InMemory {
					in, _ := inmemory.NewInMemory(zap.NewNop().Sugar())
					err := in.SetShortURL(context.Background(),
						"short123",
						"https://example.com",
						"1",
						false)
					if err != nil {
						return nil
					}
					return in
				}(), // инициализируем InMemory
				zapLog:   zap.NewNop().Sugar(),
				filePath: "/path/to/storage",
			},
			args: args{
				idShortPath: "123short",
				originalURL: "https://example.com",
				userID:      "user1",
				deleted:     false,
			},
			wantErr: true, // Ожидаем ошибку из-за пустого URL
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InFile{
				InMe:     tt.fields.InMe,
				zapLog:   tt.fields.zapLog,
				filePath: tt.fields.filePath,
			}
			err := m.SetShortURL(ctx, tt.args.idShortPath, tt.args.originalURL, tt.args.userID, tt.args.deleted)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetShortURL() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Дополнительно можно проверить, сохранился ли URL в памяти
			if !tt.wantErr {
				storedURL, err := m.GetShortURL(ctx, tt.args.idShortPath)
				if err != nil || storedURL != tt.args.originalURL {
					t.Errorf("Expected URL = %v, but got %v (error: %v)", tt.args.originalURL, storedURL, err)
				}
			}
		})
	}
}
