package inmemory

import "context"

// DeleteURLs устанавливает флаг удаления для множества записей.
func (m *InMemory) DeleteURLs(_ context.Context, urlIDs []string) error {
	if len(urlIDs) == 0 {
		return nil // Если список пуст, ничего не делаем
	}

	m.mu.Lock() // Блокируем доступ к хранилищу
	defer m.mu.Unlock()

	for _, id := range urlIDs {
		if record, exists := m.data[id]; exists {
			record.Deleted = true
			m.data[id] = record // Обновляем запись в хранилище
			m.zapLog.Infof("URL with ID %m marked as deleted", id)
		} else {
			m.zapLog.Warnf("URL with ID %m not found", id)
		}
	}

	return nil
}
