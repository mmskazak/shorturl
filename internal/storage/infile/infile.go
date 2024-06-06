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
	inMe     *inmemory.InMemory
	filePath string
}

func NewInFile(cfg *config.Config) (*InFile, error) {
	inm, err := inmemory.NewInMemory()
	if err != nil {
		return nil, fmt.Errorf("error creating inmemory storage: %w", err)
	}

	ms := &InFile{
		inMe:     inm,
		filePath: cfg.FileStoragePath,
	}

	if err := readFileStorage(ms, cfg); err != nil {
		return nil, fmt.Errorf("error read storage data: %w", err)
	}

	return ms, nil
}

func (m *InFile) GetShortURL(id string) (string, error) {
	return m.inMe.GetShortURL(id) //nolint:wrapcheck //ошибка обрабатывается далее
}

func (m *InFile) SetShortURL(id string, targetURL string) error {
	err := m.inMe.SetShortURL(id, targetURL)
	if err != nil {
		return fmt.Errorf("error setting short url: %w", err)
	}

	producer, err := rwstorage.NewProducer(m.filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания producer %w", err)
	}

	shData := rwstorage.ShortURLStruct{
		UUID:        strconv.Itoa(m.inMe.NumberOfEntries()),
		ShortURL:    id,
		OriginalURL: targetURL,
	}

	err = producer.WriteData(&shData)
	if err != nil {
		return fmt.Errorf("ошибка записи строки в файл %w", err)
	}
	producer.Close()
	log.Printf("Добавлени которкая ссылка %v", shData)

	return nil
}

func readFileStorage(m *InFile, cfg *config.Config) error {
	consumer, err := rwstorage.NewConsumer(cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("error read file storage %w", err)
	}

	for consumer.Reader.Scan() {
		dataOfURL, err := consumer.ReadLineInFile()
		if err != nil {
			return fmt.Errorf("consumer error read line in file: %w", err)
		}

		log.Printf("Прочитанные данные: %+v\n", dataOfURL)
		err = m.inMe.SetShortURL(dataOfURL.ShortURL, dataOfURL.OriginalURL)
		if err != nil {
			return fmt.Errorf("error setting short url: %w", err)
		}
		log.Printf("Длина мапы: %+v\n", m.inMe.NumberOfEntries())
	}
	return nil
}
