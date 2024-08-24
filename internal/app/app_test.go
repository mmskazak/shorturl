package app

import (
	"context"
	"testing"
	"time"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage"
	"mmskazak/shorturl/internal/storage/inmemory"

	"go.uber.org/zap"
)

func TestNewApp(t *testing.T) {
	ctx := context.Background()
	type args struct {
		cfg          *config.Config
		store        storage.Storage
		readTimeout  time.Duration
		writeTimeout time.Duration
		zapLog       *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
		want *App
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
			tt.args.zapLog)
		// Проверка типа через утверждение типа
		var _ *App = got
	}
}
