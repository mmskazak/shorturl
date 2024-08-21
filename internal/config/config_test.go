package config

import (
	"go.uber.org/zap/zapcore"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Success init config",
			want: &Config{
				Address:         ":8080",
				BaseHost:        "http://localhost:8080",
				SecretKey:       "secret",
				LogLevel:        "info",
				FileStoragePath: "/tmp/short-url-db.json",
				ReadTimeout:     10 * time.Second,
				WriteTimeout:    10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitConfig()
			if err != nil {
				t.Errorf("InitConfig() error = %v", err)
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("InitConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogLevel_Value(t *testing.T) {
	tests := []struct {
		name    string
		ll      LogLevel
		want    zapcore.Level
		wantErr bool
	}{
		{
			name:    "debug level",
			ll:      "debug",
			want:    zapcore.DebugLevel,
			wantErr: false,
		},
		{
			name:    "info level",
			ll:      "info",
			want:    zapcore.InfoLevel,
			wantErr: false,
		},
		{
			name:    "warn level",
			ll:      "warn",
			want:    zapcore.WarnLevel,
			wantErr: false,
		},
		{
			name:    "error level",
			ll:      "error",
			want:    zapcore.ErrorLevel,
			wantErr: false,
		},
		{
			name:    "dpanic level",
			ll:      "dpanic",
			want:    zapcore.DPanicLevel,
			wantErr: false,
		},
		{
			name:    "panic level",
			ll:      "panic",
			want:    zapcore.PanicLevel,
			wantErr: false,
		},
		{
			name:    "fatal level",
			ll:      "fatal",
			want:    zapcore.FatalLevel,
			wantErr: false,
		},
		{
			name:    "unknown level",
			ll:      "unknown",
			want:    zapcore.DebugLevel,
			wantErr: true,
		},
		{
			name:    "empty level",
			ll:      "",
			want:    zapcore.DebugLevel,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ll.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
		})
	}
}
