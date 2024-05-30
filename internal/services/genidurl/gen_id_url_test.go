package genidurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateShortURL(t *testing.T) {
	service := NewGenIDService()
	tests := []struct {
		name   string
		length int
		err    bool
	}{
		{
			name:   "length short URl equal 8",
			length: 8,
			err:    false,
		},
		{
			name:   "length short URl equal 5",
			length: 5,
			err:    false,
		},
		{
			name:   "length short URl equal 3",
			length: 3,
			err:    true,
		},
		{
			name:   "length short URl equal 1",
			length: 1,
			err:    true,
		},
		{
			name:   "length short URl equal 100",
			length: 100,
			err:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Generate(tt.length)
			if tt.err {
				assert.Empty(t, got)
				assert.Error(t, err)
			}

			if !tt.err {
				require.NoError(t, err)
				assert.NotEmpty(t, got)
				assert.Equal(t, tt.length, len(got), "Сгенерированный короткий URL-адрес "+
					"имеет не такую длину, как ожидалось")
			}
		})
	}
}
