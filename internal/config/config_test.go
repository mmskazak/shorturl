package config

import (
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
				Address:      ":8080",
				BaseHost:     "http://localhost:8080",
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
				LogLevel:     "info",
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