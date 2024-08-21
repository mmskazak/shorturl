package infile

import (
	"context"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/config"
	"reflect"
	"testing"
)

func TestNewInFile(t *testing.T) {
	type args struct {
		ctx    context.Context
		cfg    *config.Config
		zapLog *zap.SugaredLogger
	}
	tests := []struct {
		name    string
		args    args
		want    *InFile
		wantErr bool
	}{
		{
			name: "error due to invalid config",
			args: args{
				ctx:    context.Background(),
				cfg:    &config.Config{ /* Заполните некорректными данными */ },
				zapLog: zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInFile(tt.args.ctx, tt.args.cfg, tt.args.zapLog)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
