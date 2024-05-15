package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestInit(t *testing.T) {
	type args struct {
		level zapcore.Level
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestInit",
			args: args{
				level: zapcore.DebugLevel,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init(tt.args.level)
			require.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}
