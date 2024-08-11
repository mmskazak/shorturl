package genidurl

import "testing"

// BenchmarkGenerate function Generate()
func BenchmarkGenerate(b *testing.B) {
	genID := &GenID{}
	for i := 0; i < b.N; i++ {
		_, err := genID.Generate()
		if err != nil {
			b.Fatalf("Ошибка при генерации ID: %v", err)
		}
	}
}
