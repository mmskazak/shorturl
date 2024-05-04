package shorturl

import (
	"context"
	"mmskazak/shorturl/internal/app/config"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestApp_Start(t *testing.T) {
	t.Skip()
	type fields struct {
		config *config.Config
		router *chi.Mux
		server *http.Server
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success start App",
			fields: fields{
				config: config.InitConfig(),
				router: chi.NewRouter(),
				server: &http.Server{
					Addr:         "localhost:8080",
					Handler:      chi.NewRouter(),
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				server: tt.fields.server,
			}
			// Остановка сервера после выполнения теста
			defer func(server *http.Server, ctx context.Context) {
				err := server.Shutdown(ctx)
				require.NoError(t, err)
			}(a.server, context.Background())

			if err := a.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewApp(t *testing.T) {
	type args struct {
		cfg          *config.Config
		readTimeout  time.Duration
		writeTimeout time.Duration
	}

	initCfg := config.InitConfig()

	tests := []struct {
		name string
		args args
		want *App
	}{
		{
			name: "SuccessCreateFirstNewApp",
			args: args{
				cfg:          initCfg,
				readTimeout:  10 * time.Second,
				writeTimeout: 10 * time.Second,
			},
		},
		{
			name: "SuccessCreateSecondNewApp",
			args: args{
				cfg:          initCfg,
				readTimeout:  20 * time.Second,
				writeTimeout: 20 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewApp(tt.args.cfg, tt.args.readTimeout, tt.args.writeTimeout)
			if gotType := reflect.TypeOf(got); gotType != reflect.TypeOf(tt.want) {
				t.Errorf("NewApp() = %v, want %v", got, tt.want)
			}
		})
	}
}
