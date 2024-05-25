package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestInitWriteToOutput(t *testing.T) {
	type args struct {
		level zapcore.Level
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success InitWriteToFile",
			args: args{
				level: zapcore.DebugLevel,
			},
			wantErr: assert.NoError,
		},
		{
			name: "info level InitWriteToFile",
			args: args{
				level: zapcore.InfoLevel,
			},
			wantErr: assert.NoError,
		},
		{
			name: "warn level InitWriteToFile",
			args: args{
				level: zapcore.WarnLevel,
			},
			wantErr: assert.NoError,
		},
		{
			name: "error level InitWriteToFile",
			args: args{
				level: zapcore.ErrorLevel,
			},
			wantErr: assert.NoError,
		},
		{
			name: "invalid level InitWriteToFile",
			args: args{
				level: zapcore.Level(-1),
			},
			wantErr: assert.Error,
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
