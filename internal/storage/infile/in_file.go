package infile

import (
	"context"
	"errors"
	"fmt"
	"io"

	"go.uber.org/zap"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/services/rwstorage"
	"mmskazak/shorturl/internal/storage/inmemory"
)

const errMsgSaveBatchAndRemove = "error save batch and removing temp file %w"

// FileRecord - структура для сериализации и десериализации данных в файл.
type FileRecord struct {
	ID   string             `json:"id"`   // Короткий идентификатор URL
	Data inmemory.URLRecord `json:"data"` // Связанная информация о URL
}

type InFile struct {
	InMe     *inmemory.InMemory
	zapLog   *zap.SugaredLogger
	filePath string
}

// NewInFile - конструктор для создания нового хранилища с поддержкой работы с файлом.
func NewInFile(ctx context.Context, cfg *config.Config, zapLog *zap.SugaredLogger) (*InFile, error) {
	inm, err := inmemory.NewInMemory(zapLog)
	if err != nil {
		return nil, fmt.Errorf("error creating inmemory storage: %w", err)
	}

	ms := &InFile{
		InMe:     inm,
		filePath: cfg.FileStoragePath,
		zapLog:   zapLog,
	}

	if err := ms.readFileStorage(ctx); err != nil {
		return nil, fmt.Errorf("error read storage data: %w", err)
	}

	return ms, nil
}

// appendToFile - добавление новой записи в файл.
func (m *InFile) appendToFile(record rwstorage.ShortURLStruct) error {
	producer, err := rwstorage.NewProducer(m.filePath)
	if err != nil {
		return fmt.Errorf("error creating producer: %w", err)
	}
	defer producer.Close()

	err = producer.WriteData(&record)
	if err != nil {
		return fmt.Errorf("error writing data to file: %w", err)
	}

	m.zapLog.Infof("Added short URL: %v", record)

	return nil
}

// readFileStorage читает данные из файла и загружает их в память.
func (m *InFile) readFileStorage(ctx context.Context) error {
	consumer, err := rwstorage.NewConsumer(m.filePath)
	if err != nil {
		return fmt.Errorf("error initializing consumer: %w", err)
	}
	defer consumer.Close()

	for {
		record, err := consumer.ReadLineInFile()
		if err != nil {
			if errors.Is(err, io.EOF) || err.Error() == rwstorage.ErrEmptyFile {
				break // Конец файла достигнут
			}
			return fmt.Errorf("error reading line from file: %w", err)
		}

		if err := m.InMe.SetShortURL(ctx, record.ShortURL, record.OriginalURL, record.UserID); err != nil {
			return fmt.Errorf("error setting short URL in memory: %w", err)
		}
	}

	return nil
}

// Close - закрытие хранилища (заглушка для будущих изменений).
func (m *InFile) Close() error {
	m.zapLog.Debugln("InFile storage closed (nothing to close currently)")
	return nil
}
