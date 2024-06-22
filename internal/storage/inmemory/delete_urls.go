package inmemory

import (
	"context"
)

// DeleteURLs устанавливает флаг удаления для множества записей.
func (m *InMemory) DeleteURLs(ctx context.Context, urlIDs []string) error {
	if len(urlIDs) == 0 {
		return nil // Если список пуст, ничего не делаем
	}

	batchSize := 5000
	batches := splitSlice(urlIDs, batchSize)
	for _, batchIDs := range batches {
		err := m.deleteURLsBatch(ctx, batchIDs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *InMemory) deleteURLsBatch(_ context.Context, urlIDs []string) error {
	go func(urlIDs []string) {
		for _, id := range urlIDs {
			m.mu.Lock() // Блокируем доступ к хранилищу
			if record, exists := m.data[id]; exists {
				record.Deleted = true
				m.data[id] = record // Обновляем запись в хранилище
				m.zapLog.Infof("URL with ID %v marked as deleted", id)
			} else {
				m.zapLog.Warnf("URL with ID %v not found", id)
			}
			m.mu.Unlock()
		}
	}(urlIDs)
	return nil
}

// splitSlice разделяет исходный слайс на слайс слайсов, каждый из которых содержит не более chunkSize элементов.
func splitSlice(input []string, chunkSize int) [][]string {
	if chunkSize <= 0 {
		return nil // Если chunkSize некорректный, возвращаем nil
	}

	if len(input) == 0 {
		return [][]string{} // Если input пустой, возвращаем пустой слайс слайсов
	}

	// Рассчитываем количество частей и инициализируем result с нужным capacity
	numChunks := (len(input) + chunkSize - 1) / chunkSize
	result := make([][]string, 0, numChunks)

	// Разделяем input на части
	for i := 0; i < len(input); i += chunkSize {
		end := i + chunkSize
		if end > len(input) {
			end = len(input) // Корректируем end, чтобы не выйти за пределы слайса
		}
		result = append(result, input[i:end])
	}

	return result
}
