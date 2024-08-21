package middleware

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_compareHMAC(t *testing.T) {
	type args struct {
		sig1 string
		sig2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "identical HMACs",
			args: args{
				sig1: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue")),
				sig2: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue")),
			},
			want: true,
		},
		{
			name: "different HMACs",
			args: args{
				sig1: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue1")),
				sig2: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue2")),
			},
			want: false,
		},
		{
			name: "invalid base64 encoding for sig1",
			args: args{
				sig1: "invalidBase64",
				sig2: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue")),
			},
			want: false,
		},
		{
			name: "invalid base64 encoding for sig2",
			args: args{
				sig1: base64.RawURLEncoding.EncodeToString([]byte("testHMACValue")),
				sig2: "invalidBase64",
			},
			want: false,
		},
		{
			name: "both HMACs are invalid",
			args: args{
				sig1: "invalidBase64",
				sig2: "invalidBase64",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, compareHMAC(tt.args.sig1, tt.args.sig2), "compareHMAC(%v, %v)", tt.args.sig1, tt.args.sig2)
		})
	}
}
