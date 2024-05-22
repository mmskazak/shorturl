package mapstorage

import (
	"bufio"
	"encoding/json"
	"errors"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/services/rwstorage"
	"os"
	"path/filepath"
	"sync"
)

var ErrNotFound = errors.New("key not found")
var ErrKeyAlreadyExists = errors.New("key already exists")

type MapStorage struct {
	mu          *sync.Mutex
	Data        map[string]string
	StoragePath string
}

func NewMapStorage(pathFileStorage string) (*MapStorage, error) {
	storage := &MapStorage{
		mu:          &sync.Mutex{},
		Data:        make(map[string]string),
		StoragePath: pathFileStorage,
	}

	// Проверка, существует ли файл
	if _, err := os.Stat(pathFileStorage); os.IsNotExist(err) {
		logger.Log.Info("файла не существует")
		// Если файл не существует, создаем его, включая необходимые директории
		err = os.MkdirAll(filepath.Dir(pathFileStorage), 0750)
		if err != nil {
			logger.Log.Error("не удалось создать директории", err)
			return nil, err
		}
		logger.Log.Info("создали папки")

		file, err := os.Create(pathFileStorage)
		if err != nil {
			logger.Log.Error("не удалось создать файл", err)
			return nil, err
		}
		logger.Log.Infof("создали файл %v", pathFileStorage)
		file.Close()
	} else {
		logger.Log.Infof("файл существует %v", pathFileStorage)
		// Если файл существует, читаем его построчно
		file, err := os.Open(pathFileStorage)
		if err != nil {
			logger.Log.Error("не удалось открыть файл", err)
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			var record rwstorage.RecordURL
			err := json.Unmarshal([]byte(line), &record)
			if err != nil {
				logger.Log.Error("не удалось разобрать JSON строку", err)
				return nil, err
			}
			storage.Data[record.ShortURL] = record.OriginalURL
		}

		if err := scanner.Err(); err != nil {
			logger.Log.Error("ошибка при чтении файла", err)
			return nil, err
		}
	}

	return storage, nil
}

func (m *MapStorage) GetShortURL(id string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	targetURL, ok := m.Data[id]
	if !ok {
		return "", ErrNotFound
	}
	return targetURL, nil
}

func (m *MapStorage) SetShortURL(id string, targetURL string) error {
	if id == "" {
		return errors.New("id is empty")
	}
	if targetURL == "" {
		return errors.New("URL is empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Data[id]; ok {
		return ErrKeyAlreadyExists
	}
	m.Data[id] = targetURL

	// Открыть файл в режиме добавления
	file, err := os.OpenFile(m.StoragePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Создать новую запись
	record := rwstorage.RecordURL{
		ShortURL:    id,
		OriginalURL: targetURL,
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Записать новую запись в файл
	if _, err := file.WriteString(string(recordJSON) + "\n"); err != nil {
		return err
	}

	return nil
}
