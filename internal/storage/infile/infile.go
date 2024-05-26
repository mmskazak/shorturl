package infile

import (
	"fmt"
	"log"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/services/rwstorage"
	"mmskazak/shorturl/internal/storage/inmemory"
	"strconv"
)

type InFile struct {
	InMe     *inmemory.InMemory
	FilePath string
}

func NewInFile(cfg *config.Config) (*InFile, error) {
	inm, err := inmemory.NewInMemory()
	if err != nil {
		return nil, fmt.Errorf("error creating inmemory storage: %w", err)
	}

	ms := &InFile{
		InMe:     inm,
		FilePath: cfg.FileStoragePath,
	}

	if err := readFileStorage(ms, cfg); err != nil {
		return nil, fmt.Errorf("error read storage data: %w", err)
	}

	return ms, nil
}

func (m *InFile) GetShortURL(id string) (string, error) {
	return m.InMe.GetShortURL(id) //nolint:wrapcheck //ошибка обрабатывается далее
}

func (m *InFile) SetShortURL(id string, targetURL string) error {
	err := m.InMe.SetShortURL(id, targetURL)
	if err != nil {
		return fmt.Errorf("error setting short url: %w", err)
	}

	if m.FilePath != "" {
		producer, err := rwstorage.NewProducer(m.FilePath)
		if err != nil {
			return fmt.Errorf("ошибка создания producer %w", err)
		}

		shData := rwstorage.ShortURLStruct{
			UUID:        strconv.Itoa(len(m.InMe.Data)),
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
		m.InMe.Data[dataOfURL.ShortURL] = dataOfURL.OriginalURL
		log.Printf("Длина мапы: %+v\n", len(m.InMe.Data))
	}
	return nil
}
