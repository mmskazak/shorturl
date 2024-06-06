package infile

import (
	"fmt"
	"log"
	"mmskazak/shorturl/internal/services/rwstorage"
	"mmskazak/shorturl/internal/storage"
	"os"
	"strconv"
)

func (m *InFile) SaveBatch(items []storage.Incoming, baseHost string) ([]storage.Output, error) {
	lenItems := len(items)
	outputs, err := m.inMe.SaveBatch(items, baseHost)
	if err != nil {
		return nil, fmt.Errorf("error saving batch infile: %w", err)
	}

	tempFilePath := m.filePath + ".tmp"
	producer, err := rwstorage.NewProducer(tempFilePath)
	if err != nil {
		return nil, fmt.Errorf("save batch ошибка создания producer %w", err)
	}
	defer producer.Close()

	batchSize := 4
	number := m.inMe.NumberOfEntries() - lenItems
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
				err := os.Remove(tempFilePath)
				if err != nil {
					return nil, fmt.Errorf(errMsgSaveBatchAndRemove, err)
				}
				return nil, fmt.Errorf("save batch ошибка записи строки в файл %w", err)
			}
			batch = batch[:0] // Очистить батч
		}
	}

	// Записать оставшиеся данные, если они есть
	if len(batch) > 0 {
		err = producer.WriteBatch(batch)
		if err != nil {
			err := os.Remove(tempFilePath)
			if err != nil {
				return nil, fmt.Errorf(errMsgSaveBatchAndRemove, err)
			}
			return nil, fmt.Errorf("save batch ошибка записи строки в файл %w", err)
		}
	}

	// Переносим логику дописывания данных в основной файл в Producer
	producerForCopy, err := rwstorage.NewProducer(m.filePath)
	if err != nil {
		err := os.Remove(tempFilePath)
		if err != nil {
			return nil, fmt.Errorf(errMsgSaveBatchAndRemove, err)
		}
		return nil, fmt.Errorf("save batch ошибка создания producer для основного файла %w", err)
	}
	defer producerForCopy.Close()

	err = producer.AppendToFile(tempFilePath, m.filePath)
	if err != nil {
		err := os.Remove(tempFilePath)
		if err != nil {
			return nil, fmt.Errorf(errMsgSaveBatchAndRemove, err)
		}
		return nil, fmt.Errorf("save batch ошибка копирования данных из временного файла %w", err)
	}

	err = os.Remove(tempFilePath)
	if err != nil {
		log.Printf("os remove temp file %s error: %v", tempFilePath, err)
	}

	return outputs, nil
}
