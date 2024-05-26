package infile

import (
	"errors"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/services/rwstorage"
	"mmskazak/shorturl/internal/storage/inmemory"
	"strconv"
)

var ErrNotFound = errors.New("key not found")
var ErrKeyAlreadyExists = errors.New("key already exists")

type InFile struct {
	*inmemory.InMemory
	FilePath string
}

func NewInFile(cfg *config.Config) (*InFile, error) {
	inm, err := inmemory.NewInMemory()
	if err != nil {
		return nil, fmt.Errorf("error creating inmemory storage: %w", err)
	}

	ms := &InFile{
		InMemory: inm,
		FilePath: cfg.FileStoragePath,
	}

	if err := readFileStorage(ms, cfg); err != nil {
		return nil, fmt.Errorf("error read storage data: %w", err)
	}

	return ms, nil
}

func (m *InFile) GetShortURL(id string) (string, error) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	targetURL, ok := m.Data[id]
	if !ok {
		return "", ErrNotFound
	}
	return targetURL, nil
}

func (m *InFile) SetShortURL(id string, targetURL string) error {
	if id == "" {
		return errors.New("id is empty")
	}
	if targetURL == "" {
		return errors.New("URL is empty")
	}
	m.Mu.Lock()
	defer m.Mu.Unlock()
	if _, ok := m.Data[id]; ok {
		return ErrKeyAlreadyExists
	}
	m.Data[id] = targetURL

	if m.FilePath != "" {
		producer, err := rwstorage.NewProducer(m.FilePath)
		if err != nil {
			return fmt.Errorf("ошибка создания producer %w", err)
		}

		shData := rwstorage.ShortURLStruct{
			UUID:        strconv.Itoa(len(m.Data)),
			ShortURL:    id,
			OriginalURL: targetURL,
		}

		err = producer.WriteData(&shData)
		if err != nil {
			return fmt.Errorf("ошибка записи строки в файл %w", err)
		}
		producer.Close()
		log.Printf("Добавлени которкая ссылка %v", shData)
	}
	return nil
}

func readFileStorage(m *InFile, cfg *config.Config) error {
	consumer, err := rwstorage.NewConsumer(cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("error read file storage %w", err)
	}

	for {
		dataOfURL, err := consumer.ReadDataFromFile()
		if err != nil {
			if err.Error() != "EOF" {
				return fmt.Errorf("ошибка при чтении: %w", err)
			}
			fmt.Println("Достигнут конец файла.")
			break
		}

		log.Printf("Прочитанные данные: %+v\n", dataOfURL)
		m.Data[dataOfURL.ShortURL] = dataOfURL.OriginalURL
		log.Printf("Длина мапы: %+v\n", len(m.Data))
	}
	return nil
}
