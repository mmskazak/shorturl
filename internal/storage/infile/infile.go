package infile

import (
	"context"
	"fmt"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/services/rwstorage"
	"mmskazak/shorturl/internal/storage/inmemory"
	"strconv"

	"go.uber.org/zap"
)

// Определяем константу для сообщения об ошибке.
const errMsgSaveBatchAndRemove = "error save batch and removing temp file %w"

type InFile struct {
	inMe     *inmemory.InMemory
	zapLog   *zap.SugaredLogger
	filePath string
}

func NewInFile(ctx context.Context, cfg *config.Config, zapLog *zap.SugaredLogger) (*InFile, error) {
	inm, err := inmemory.NewInMemory(zapLog)
	if err != nil {
		return nil, fmt.Errorf("error creating inmemory storage: %w", err)
	}

	ms := &InFile{
		inMe:     inm,
		filePath: cfg.FileStoragePath,
		zapLog:   zapLog,
	}

	if err := readFileStorage(ctx, ms, cfg); err != nil {
		return nil, fmt.Errorf("error read storage data: %w", err)
	}

	return ms, nil
}

func (m *InFile) GetShortURL(ctx context.Context, id string) (string, error) {
	return m.inMe.GetShortURL(ctx, id) //nolint:wrapcheck //ошибка обрабатывается далее
}

func (m *InFile) SetShortURL(ctx context.Context, id string, targetURL string) error {
	err := m.inMe.SetShortURL(ctx, id, targetURL)
	if err != nil {
		return fmt.Errorf("error set short url: %w", err)
	}

	producer, err := rwstorage.NewProducer(m.filePath)
	if err != nil {
		return fmt.Errorf("erorr create producer %w", err)
	}

	shData := rwstorage.ShortURLStruct{
		UUID:        strconv.Itoa(m.inMe.NumberOfEntries()),
		ShortURL:    id,
		OriginalURL: targetURL,
	}

	err = producer.WriteData(&shData)
	if err != nil {
		return fmt.Errorf("error write string in file %w", err)
	}
	producer.Close()
	m.zapLog.Infof("Add short link %v", shData)

	return nil
}

func readFileStorage(ctx context.Context, m *InFile, cfg *config.Config) error {
	consumer, err := rwstorage.NewConsumer(cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("error read file storage %w", err)
	}

	for consumer.Reader.Scan() {
		dataOfURL, err := consumer.ReadLineInFile()
		if err != nil {
			return fmt.Errorf("consumer error read line in file: %w", err)
		}

		m.zapLog.Infof("Readed data: %+v\n", dataOfURL)
		err = m.inMe.SetShortURL(ctx, dataOfURL.ShortURL, dataOfURL.OriginalURL)
		if err != nil {
			return fmt.Errorf("error setting short url: %w", err)
		}
		m.zapLog.Infof("Lenght map storage: %+v\n", m.inMe.NumberOfEntries())
	}
	return nil
}

func (m *InFile) Close() error {
	// На данный момент закрывать нечего, но метод оставлен для возможных будущих изменений
	m.zapLog.Debugln("InFile storage closed (nothing to close currently)")
	return nil
}
