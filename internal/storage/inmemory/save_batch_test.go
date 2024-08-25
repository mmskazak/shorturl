package inmemory

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/contracts/mocks"
	"mmskazak/shorturl/internal/models"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestInMemory_SaveBatch(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctxBg := context.Background()

	mockGenID := mocks.NewMockIGenIDForURL(ctrl)

	output := []models.Output{
		{
			CorrelationID: "123",
			ShortURL:      "http://127.0.0.1:8080/QwErAsDf",
		},
		{
			CorrelationID: "456",
			ShortURL:      "http://127.0.0.1:8080/AnOtHeR",
		},
	}

	incoming := []models.Incoming{
		{
			CorrelationID: "123",
			OriginalURL:   "https://example.com/long-url-00012",
		},
		{
			CorrelationID: "456",
			OriginalURL:   "https://example.com/long-url-00013",
		},
	}

	type fields struct {
		mu        *sync.Mutex
		data      map[string]models.URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	type args struct {
		items     []models.Incoming
		baseHost  string
		userID    string
		generator contracts.IGenIDForURL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Output
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]models.URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				items:     []models.Incoming{},
				baseHost:  "http://127.0.0.1:8080",
				userID:    "1",
				generator: mockGenID,
			},
			want:    []models.Output{},
			wantErr: false,
		},
		{
			name: "test 2",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]models.URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				items:     incoming,
				baseHost:  "http://127.0.0.1:8080",
				userID:    "1",
				generator: mockGenID,
			},
			want:    output,
			wantErr: false,
		},
		{
			name: "test 3",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]models.URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(),
			},
			args: args{
				items:     incoming,
				baseHost:  "http://127.0.0.1:8080",
				userID:    "1",
				generator: mockGenID,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "test 2" {
				mockGenID.EXPECT().Generate().Return("QwErAsDf", nil).Times(1)
				mockGenID.EXPECT().Generate().Return("AnOtHeR", nil).Times(1)
			}

			if tt.name == "test 3" {
				mockGenID.EXPECT().Generate().Return("QwErWeRt", nil).AnyTimes()
			}

			m := &InMemory{
				mu:        tt.fields.mu,
				data:      tt.fields.data,
				userIndex: tt.fields.userIndex,
				zapLog:    tt.fields.zapLog,
			}
			got, err := m.SaveBatch(ctxBg, tt.args.items, tt.args.baseHost, tt.args.userID, tt.args.generator)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveBatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
