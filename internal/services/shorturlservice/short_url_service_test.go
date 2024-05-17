package shorturlservice

import (
	"fmt"
	"mmskazak/shorturl/internal/storage/mapstorage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testID = "TeSt0001"

type GenerateIDDummy struct{}

func (g *GenerateIDDummy) Generate(_ int) (string, error) {
	return testID, nil
}

func TestShortURLService_GenerateShortURL(t *testing.T) {
	type args struct {
		dto       DTOShortURL
		generator IGenIDForURL
		storage   IStorage
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				dto: DTOShortURL{
					OriginalURL:  "http://ya.ru",
					BaseHost:     "http://localhost.com",
					MaxIteration: 10,
					LengthID:     8,
				},
				generator: &GenerateIDDummy{},
				storage:   mapstorage.NewMapStorage(),
			},
			want:    "http://localhost.com/TeSt0001",
			wantErr: assert.NoError,
		},
		{
			name: "original url is empty",
			args: args{
				dto: DTOShortURL{
					OriginalURL:  "",
					BaseHost:     "http://localhost.com",
					MaxIteration: 10,
					LengthID:     8,
				},
				generator: &GenerateIDDummy{},
				storage:   mapstorage.NewMapStorage(),
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "base host is empty",
			args: args{
				dto: DTOShortURL{
					OriginalURL:  "http://ya.ru",
					BaseHost:     "",
					MaxIteration: 10,
					LengthID:     10,
				},
				generator: &GenerateIDDummy{},
				storage:   mapstorage.NewMapStorage(),
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "length id is small",
			args: args{
				dto: DTOShortURL{
					OriginalURL:  "http://ya.ru",
					BaseHost:     "http://localhost.com",
					MaxIteration: 5,
					LengthID:     5,
				},
				generator: &GenerateIDDummy{},
				storage: func() *mapstorage.MapStorage {
					msTest := mapstorage.NewMapStorage()
					err := msTest.SetShortURL(testID, "http://google.com")
					if err != nil {
						require.NoError(t, err)
					}
					return msTest
				}(),
			},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortURLService{}
			got, err := s.GenerateShortURL(tt.args.dto, tt.args.generator, tt.args.storage)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateShortURL(%v, %v, %v)",
				tt.args.dto, tt.args.generator, tt.args.storage)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GenerateShortURL(%v, %v, %v)",
				tt.args.dto, tt.args.generator, tt.args.storage)
		})
	}
}
