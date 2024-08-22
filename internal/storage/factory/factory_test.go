package factory

import (
	"context"
	"fmt"
	"testing"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage"
	"mmskazak/shorturl/internal/storage/inmemory"

	"go.uber.org/zap"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
		want    storage.Storage
	}{
		{
			name: "InMemory storage",
			cfg: &config.Config{
				DataBaseDSN:     "",
				FileStoragePath: "",
			},
			wantErr: false,
			want:    &inmemory.InMemory{}, // Убедитесь, что этот тип соответствует вашему типу для InMemory
		},
		{
			name: "Error on PostgreSQL initialization",
			cfg: &config.Config{
				DataBaseDSN: "invalid_dsn",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Error on InFile initialization",
			cfg: &config.Config{
				FileStoragePath: "/invalid/path/to/storage/file.db",
			},
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zapLog := zap.NewExample().Sugar() // Создание логгера для тестов
			got, err := NewStorage(context.Background(), tt.cfg, zapLog)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Проверяем, что тип возвращаемого хранилища соответствует ожидаемому
			if tt.want != nil && got != nil && fmt.Sprintf("%T", got) != fmt.Sprintf("%T", tt.want) {
				t.Errorf("NewStorage() got = %T, want %T", got, tt.want)
			}
		})
	}
}
