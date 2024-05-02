package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortURL(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "Success generate short url",
			length: 8,
		},
		{
			name:   "Success generate short url",
			length: 5,
		},
		{
			name:   "Success generate short url",
			length: 1,
		},
		{
			name:   "Success generate short url",
			length: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateShortURL(tt.length)
			assert.NotEmpty(t, got)
			// Проверяем, что got является строкой
			assert.IsType(t, "", got, "GenerateShortURL() should return a string")
			// Проверяем длину полученной строки
			assert.Equal(t, tt.length, len(got), "Generated short URL length is not as expected") //nolint:testifylint
		})
	}
}
