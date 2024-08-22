package infile

import (
	"context"
	"reflect"
	"testing"
	"time"

	"mmskazak/shorturl/internal/config"

	"go.uber.org/zap"
)

func TestNewInFile(t *testing.T) {
	ctx := context.Background()

	type args struct {
		cfg    *config.Config
		zapLog *zap.SugaredLogger
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "error due to invalid config",
			args: args{
				cfg:    &config.Config{},
				zapLog: zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "error due to invalid config",
			args: args{

				cfg: &config.Config{
					Address:         ":8080",
					BaseHost:        "http://localhost:8080",
					SecretKey:       "secret",
					LogLevel:        "info",
					FileStoragePath: "/tmp/short-url-db.json",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
				},
				zapLog: zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
			},
			want:    "/tmp/short-url-db.json",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInFile(ctx, tt.args.cfg, tt.args.zapLog)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr {
				return
			}

			// Используем рефлексию для доступа к приватному полю
			r := reflect.ValueOf(got).Elem()
			field := r.FieldByName("filePath")
			if field.String() != tt.want {
				t.Errorf("NewInFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
