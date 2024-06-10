package infile

import (
	"context"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"mmskazak/shorturl/internal/storage/inmemory"
)

func createTempFile(t *testing.T, content string) string {
	t.Helper() // Mark this function as a test helper

	tmpFile, err := os.CreateTemp("", "shorturl-*.json")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}

	err = tmpFile.Close()

	if err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}

func TestInFile_GetShortURL(t *testing.T) {
	// Логгер
	logger, err := zap.NewProduction()
	require.NoError(t, err)

	// Не забываем освобождать ресурсы после завершения теста
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Warn("Failed to sync logger", zap.Error(err))
		}
	}(logger)

	type fields struct {
		InMe     *inmemory.InMemory
		FilePath string
	}

	type args struct {
		id string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success has data",
			fields: fields{
				InMe: func() *inmemory.InMemory {
					inm, err := inmemory.NewInMemory(logger.Sugar())
					require.NoError(t, err)
					ctx := context.TODO()
					err = inm.SetShortURL(ctx, "test0001", "https://google.com")
					require.NoError(t, err)
					return inm
				}(),
				FilePath: createTempFile(t, ""),
			},
			args: args{
				id: "test0001",
			},
			want:    "https://google.com",
			wantErr: false,
		},
		{
			name: "error not found",
			fields: fields{
				InMe: func() *inmemory.InMemory {
					inm, err := inmemory.NewInMemory(logger.Sugar())
					require.NoError(t, err)
					ctx := context.TODO()
					err = inm.SetShortURL(ctx, "test0001", "https://google.com")
					require.NoError(t, err)
					return inm
				}(),
				FilePath: createTempFile(t, ""),
			},
			args: args{
				id: "test0002",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "error not found",
			fields: fields{
				InMe: func() *inmemory.InMemory {
					inm, err := inmemory.NewInMemory(logger.Sugar())
					require.NoError(t, err)
					ctx := context.TODO()
					err = inm.SetShortURL(ctx, "test0001", "https://google.com")
					require.NoError(t, err)
					return inm
				}(),
				FilePath: createTempFile(t, ""),
			},
			args: args{
				id: "test0",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InFile{
				inMe:     tt.fields.InMe,
				filePath: tt.fields.FilePath,
				zapLog:   logger.Sugar(),
			}
			ctx := context.TODO()
			got, err := m.GetShortURL(ctx, tt.args.id)
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
