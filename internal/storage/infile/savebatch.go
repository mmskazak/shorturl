package infile

import (
	"context"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/services/rwstorage"
	"mmskazak/shorturl/internal/storage"
	"os"
	"strconv"
)

func (m *InFile) SaveBatch(
	ctx context.Context,
	items []storage.Incoming,
	baseHost string,
	userID string,
	generator storage.IGenIDForURL,
) ([]storage.Output, error) {
	lenItems := len(items)
	outputs, err := m.inMe.SaveBatch(ctx, items, baseHost, userID, generator)
	if err != nil {
		return nil, err // прокидываем оригинальную ошибку
	}

	tempFilePath := m.filePath + ".tmp"
	producer, err := rwstorage.NewProducer(tempFilePath)
	if err != nil {
		return nil, fmt.Errorf("save batch ошибка создания producer %w", err)
	}
	defer producer.Close()

	batchSize := 1000
	number := m.inMe.NumberOfEntries() - lenItems
	batch := make([]rwstorage.ShortURLStruct, 0, batchSize)

	for _, item := range items {
		number++
		shortURLStruct := rwstorage.ShortURLStruct{
			ID:          strconv.Itoa(number),
			ShortURL:    item.CorrelationID,
			OriginalURL: item.OriginalURL,
			UserID:      userID,
			Deleted:     false,
		}

		batch = append(batch, shortURLStruct)
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
			return nil, fmt.Errorf("save batch error write lint in file: %w", err)
		}
	}

	err = producer.AppendToFile(tempFilePath, m.filePath)
	if err != nil {
		err := os.Remove(tempFilePath)
		if err != nil {
			return nil, fmt.Errorf(errMsgSaveBatchAndRemove, err)
		}
		return nil, fmt.Errorf("save batch error copy from temp to main file: %w", err)
	}

	err = os.Remove(tempFilePath)
	if err != nil {
		log.Printf("os remove temp file %s error: %v", tempFilePath, err)
	}

	return outputs, nil
}
