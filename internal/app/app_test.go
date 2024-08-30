package app

import (
	"context"
	"mmskazak/shorturl/internal/contracts"
	"testing"
	"time"

	"mmskazak/shorturl/internal/services/shorturlservice"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage/inmemory"

	"go.uber.org/zap"
)

func TestNewApp(t *testing.T) {
	ctx := context.Background()
	type args struct {
		store           contracts.Storage          // Зависит от реализации
		shortURLService contracts.IShortURLService // Зависит от реализации
		zapLog          *zap.SugaredLogger         // 8 байт
		cfg             *config.Config             // 8 байт
		readTimeout     time.Duration              // 8 байт
		writeTimeout    time.Duration              // 8 байт
	}
	tests := []struct {
		want *App
		name string
		args args
	}{
		{
			name: "test 1",
			args: args{
				cfg: &config.Config{
					Address: "https://127.0.0.1",
				},
				store: func() *inmemory.InMemory {
					in, _ := inmemory.NewInMemory(zap.NewNop().Sugar())
					return in
				}(),
				shortURLService: shorturlservice.NewShortURLService(),
			},
			want: &App{},
		},
	}
	for _, tt := range tests {
		got := NewApp(ctx,
			tt.args.cfg,
			tt.args.store,
			tt.args.readTimeout,
			tt.args.writeTimeout,
			tt.args.zapLog,
			tt.args.shortURLService,
		)
		// Проверка типа через утверждение типа
		var _ *App = got
	}
}
