package inmemory

import (
	"context"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
)

// MarkURLAsDeleted - помечает URL как удаленный.
func (m *InMemory) MarkURLAsDeleted(_ context.Context, id string) error {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	record, ok := m.Data[id]
	if !ok {
		return storageErrors.ErrNotFound
	}

	record.Deleted = true
	m.Data[id] = record // Обновляем запись

	return nil
}
