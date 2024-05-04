package config

import (
	"reflect"
	"testing"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Success init config",
			want: &Config{
				address:  ":8080",
				baseHost: "http://localhost:8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
