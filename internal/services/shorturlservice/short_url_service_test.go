package shorturlservice

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage"
	"mmskazak/shorturl/internal/storage/infile"
	"mmskazak/shorturl/internal/storage/inmemory"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testID = "TeSt0001"
)

type GenerateIDDummy struct{}

func (g *GenerateIDDummy) Generate(_ int) (string, error) {
	return testID, nil
}

func createTempFile(t *testing.T, content string) string {
	t.Helper() // Mark this function as a test helper

	tmpFile, err := os.CreateTemp("", "shorturl-*.json")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}

	err = tmpFile.Close()

	if err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}

func TestShortURLService_GenerateShortURL(t *testing.T) {
	type args struct {
		dto       DTOShortURL
		generator IGenIDForURL
		data      storage.Storage
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
				data: func() *inmemory.InMemory {
					s, err := inmemory.NewInMemory()
					require.NoError(t, err)
					return s
				}(),
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
				data: func() *inmemory.InMemory {
					s, err := inmemory.NewInMemory()
					require.NoError(t, err)
					return s
				}(),
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
				data: func() *inmemory.InMemory {
					s, err := inmemory.NewInMemory()
					require.NoError(t, err)
					return s
				}(),
			},
			want:    "",
			wantErr: assert.Error,
		},
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
				data: func() *infile.InFile {
					cfg := config.Config{
						FileStoragePath: createTempFile(t, ""),
					}
					s, err := infile.NewInFile(&cfg)
					require.NoError(t, err)
					return s
				}(),
			},
			want:    "http://localhost.com/TeSt0001",
			wantErr: assert.NoError,
		},
		{
			name: "success with InFile",
			args: args{
				dto: DTOShortURL{
					OriginalURL:  "http://ya.ru",
					BaseHost:     "http://localhost.com",
					MaxIteration: 10,
					LengthID:     8,
				},
				generator: &GenerateIDDummy{},
				data: func() *infile.InFile {
					cfg := config.Config{
						FileStoragePath: createTempFile(t, ""),
					}
					s, err := infile.NewInFile(&cfg)
					require.NoError(t, err)
					return s
				}(),
			},
			want:    "http://localhost.com/TeSt0001",
			wantErr: assert.NoError,
		},
		{
			name: "error with InFile",
			args: args{
				dto: DTOShortURL{
					OriginalURL:  "http://ya.ru",
					BaseHost:     "http://localhost.com",
					MaxIteration: 10,
					LengthID:     8,
				},
				generator: &GenerateIDDummy{},
				data: func() *infile.InFile {
					cfg := config.Config{
						FileStoragePath: createTempFile(t, ""),
					}
					s, err := infile.NewInFile(&cfg)
					require.NoError(t, err)

					ctx := context.TODO()
					err = s.SetShortURL(ctx, testID, "http://ya.ru")
					require.NoError(t, err)
					return s
				}(),
			},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortURLService{}
			got, err := s.GenerateShortURL(tt.args.dto, tt.args.generator, tt.args.data)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateShortURL(%v, %v, %v)",
				tt.args.dto, tt.args.generator, tt.args.data)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GenerateShortURL(%v, %v, %v)",
				tt.args.dto, tt.args.generator, tt.args.data)
		})
	}
}
