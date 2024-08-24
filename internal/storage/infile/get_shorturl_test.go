package infile

import (
	"context"
	"testing"

	"mmskazak/shorturl/internal/storage/inmemory"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInFile_GetShortURL(t *testing.T) {
	// Создание нового контроллера для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Контекст и создание моков
	ctx := context.Background()

	// Создание экземпляра InMemory
	inMemory, err := inmemory.NewInMemory(zap.NewNop().Sugar())
	if err != nil {
		t.Fatalf("Failed to create InMemory instance: %v", err)
	}

	// Добавление тестовых данных в InMemory
	testID := "testID"
	OriginalURL := "http://example.com"
	expectedURL := "http://example.com"
	err = inMemory.SetShortURL(ctx, testID, OriginalURL, "user1", false)
	require.NoError(t, err)

	// Создание экземпляра InFile
	s := &InFile{
		InMe:     inMemory,
		zapLog:   zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
		filePath: "/mock/path",
	}

	// Определение тестовых случаев
	tests := []struct {
		name string
		args struct {
			id string
		}
		want    string
		wantErr bool
	}{
		{
			name: "Success case",
			args: struct {
				id string
			}{
				id: testID,
			},
			want:    expectedURL,
			wantErr: false,
		},
		{
			name: "Error case: ID not found",
			args: struct {
				id string
			}{
				id: "nonExistingID",
			},
			want:    "",
			wantErr: true,
		},
	}

	// Запуск тестов÷
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetShortURL(ctx, tt.args.id)
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
