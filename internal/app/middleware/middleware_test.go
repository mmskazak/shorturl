package middleware //nolint: golint

import (
	"net/http"
	"reflect"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Создаем фейковый обработчик для тестирования
	fakeHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	type args struct {
		next http.Handler
	}

	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{
			name: "Test if returns http.HandlerFunc",
			args: args{
				next: fakeHandler,
			},
			want: reflect.TypeOf(http.HandlerFunc(nil)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoggingMiddleware(tt.args.next)
			if gotType := reflect.TypeOf(got); gotType != tt.want {
				t.Errorf("LoggingMiddleware() returned a handler of type %v, want %v", gotType, tt.want)
			}
		})
	}
}
