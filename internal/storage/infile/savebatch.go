package infile

import (
	"fmt"
	"mmskazak/shorturl/internal/services/rwstorage"
	"mmskazak/shorturl/internal/storage"
	"strconv"
)

func (m *InFile) SaveBatch(items []storage.Incoming, baseHost string) ([]storage.Output, error) {
	outputs, err := m.inMe.SaveBatch(items, baseHost)
	if err != nil {
		return nil, fmt.Errorf("error saving batch infile: %w", err)
	}

	producer, err := rwstorage.NewProducer(m.filePath)
	if err != nil {
		return nil, fmt.Errorf("save batch ошибка создания producer %w", err)
	}
	defer producer.Close()

	batchSize := 4
	number := m.inMe.NumberOfEntries() - len(items)
	batch := make([]rwstorage.ShortURLStruct, 0, batchSize)

	for _, item := range items {
		number++
		shData := rwstorage.ShortURLStruct{
			UUID:        strconv.Itoa(number),
			ShortURL:    item.CorrelationID,
			OriginalURL: item.OriginalURL,
		}

		batch = append(batch, shData)
		if len(batch) == batchSize {
			err = producer.WriteBatch(batch)
			if err != nil {
				return nil, fmt.Errorf("save batch ошибка записи нескольких строк в файл %w", err)
			}
			batch = batch[:0] // Очистить батч
		}
	}

	// Записать оставшиеся данные, если они есть
	if len(batch) > 0 {
		err = producer.WriteBatch(batch)
		if err != nil {
			return nil, fmt.Errorf("save batch ошибка записи строки в файл %w", err)
		}
	}

	return outputs, nil
}
