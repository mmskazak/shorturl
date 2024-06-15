package infile

import "context"

// GetShortURL - получение оригинального URL по короткому идентификатору.
func (m *InFile) GetShortURL(ctx context.Context, id string) (string, error) {
	return m.inMe.GetShortURL(ctx, id) //nolint:wrapcheck // ошибка обрабатывается далее
}
