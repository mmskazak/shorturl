package middleware

import (
	"errors"
	"net/http"
	"testing"
)

type ResponseWriterDummy struct{}

func (w *ResponseWriterDummy) Header() http.Header {
	return http.Header{}
}

func (w *ResponseWriterDummy) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, errors.New("тестовая ошибка")
	}
	return len(b), nil
}

func (w *ResponseWriterDummy) WriteHeader(int) {}

func Test_loggingResponseWriter_Write(t *testing.T) {
	rwd := ResponseWriterDummy{}

	type fields struct {
		ResponseWriter http.ResponseWriter
		responseData   *responseData
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "success case: should return a responseData",
			fields: fields{
				ResponseWriter: &rwd,
				responseData:   &responseData{},
			},
			args: args{
				b: []byte("hello world"),
			},
			want:    len("hello world"),
			wantErr: false,
		},
		{
			name: "error case: should return an error",
			fields: fields{
				ResponseWriter: &rwd,
				responseData:   &responseData{},
			},
			args: args{
				b: []byte(""),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &loggingResponseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				responseData:   tt.fields.responseData,
			}
			got, err := r.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}
