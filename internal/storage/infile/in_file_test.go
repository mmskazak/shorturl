package infile

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			name: "success create config struct",
			args: args{

				cfg: &config.Config{
					Address:         ":8080",
					BaseHost:        "http://localhost:8080",
					SecretKey:       "secret",
					LogLevel:        "info",
					FileStoragePath: "/tmp/success-create-config.json",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
				},
				zapLog: zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
			},
			want:    "/tmp/success-create-config.json",
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
				t.Errorf("NewInFile() field.String() = %v, want %v", field.String(), tt.want)
			}
		})
	}
}

func Test_parseShortURLStruct_Success(t *testing.T) {
	validJSON := `{
    "id": "1",
    "short_url": "testtest",
    "original_url": "http://original.url",
    "user_id": "user123",
    "deleted": false
	}`
	sURL := shortURLStruct{
		ID:          "1",
		ShortURL:    "testtest",
		OriginalURL: "http://original.url",
		UserID:      "user123",
		Deleted:     false,
	}

	got, err := parseShortURLStruct(validJSON)
	require.NoError(t, err)
	assert.Equal(t, sURL, got)
}

func Test_parseShortURLStruct_Err(t *testing.T) {
	validJSON := `{"id": "1",`
	sURL := shortURLStruct{}

	got, err := parseShortURLStruct(validJSON)
	assert.Error(t, err)
	assert.Equal(t, sURL, got)
}
