package inmemory

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"reflect"
	"sync"
	"testing"
)

func TestNewInMemory(t *testing.T) {
	type args struct {
		zapLog *zap.SugaredLogger
	}

	zaplogSugar := zaptest.NewLogger(t).Sugar()

	tests := []struct {
		name    string
		args    args
		want    *InMemory
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				zapLog: zaplogSugar,
			},
			want: &InMemory{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zaplogSugar,
			},
			wantErr: false,
		},
		{
			name: "test 2",
			args: args{
				zapLog: nil,
			},
			want: &InMemory{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInMemory(tt.args.zapLog)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInMemory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemory() got = %v, want %v", got, tt.want)
			}
		})
	}
}
