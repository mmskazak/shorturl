package inmemory

import (
	"context"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

// MarkURLAsDeleted - помечает URL как удаленный.
func (m *InMemory) MarkURLAsDeleted(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	record, ok := m.data[id]
	if !ok {
		return storageErrors.ErrNotFound
	}

	record.Deleted = true
	m.data[id] = record // Обновляем запись

	return nil
}
