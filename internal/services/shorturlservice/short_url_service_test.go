package shorturlservice

import (
	"context"
	"github.com/golang/mock/gomock"
	"mmskazak/shorturl/internal/services/shorturlservice/mocks"
	"testing"
)

func TestShortURLService_GenerateShortURL(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	type fields struct {
		maxIteration int
	}
	type args struct {
		dto       DTOShortURL
		generator *mocks.MockIGenIDForURL
		data      *mocks.MockISetShortURL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				maxIteration: 10,
			},
			args: args{
				dto: DTOShortURL{
					OriginalURL: "ya.ru",
					UserID:      "1",
					BaseHost:    "http://localhost",
					Deleted:     false,
				},
				generator: mocks.NewMockIGenIDForURL(ctrl),
				data:      mocks.NewMockISetShortURL(ctrl),
			},
			want:    "http://localhost/expectedID",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortURLService{
				maxIteration: tt.fields.maxIteration,
			}
			// Настройка ожиданий
			tt.args.generator.EXPECT().Generate().Return("expectedID", nil)
			tt.args.data.EXPECT().SetShortURL(ctx,
				"expectedID",
				tt.args.dto.OriginalURL,
				tt.args.dto.UserID,
				false).
				Return(nil)

			got, err := s.GenerateShortURL(ctx, tt.args.dto, tt.args.generator, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateShortURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
