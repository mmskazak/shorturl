package infile

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/storage/inmemory"
)

type InFile struct {
	InMe     *inmemory.InMemory
	zapLog   *zap.SugaredLogger
	filePath string
}

type shortURLStruct struct {
	ID          string `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	Deleted     bool   `json:"deleted"`
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

// parseShortURLStruct - вспомогательная функция для парсинга строки JSON в структуру shortURLStruct.
func parseShortURLStruct(line string) (shortURLStruct, error) {
	var record shortURLStruct
	err := json.Unmarshal([]byte(line), &record)
	if err != nil {
		return shortURLStruct{}, fmt.Errorf("error parsing JSON: %w", err)
	}
	return record, nil
}

func (m *InFile) readFileStorage(ctx context.Context) error {
	// Открываем файл
	file, err := os.Open(m.filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			m.zapLog.Warnf("error close file %w", err)
		}
	}(file)

	// Читаем файл построчно
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Парсим строку JSON в структуру shortURLStruct
		record, err := parseShortURLStruct(line)
		if err != nil {
			return fmt.Errorf("error parsing line from file: %w", err)
		}

		// Добавляем запись в InMemoryStorage
		if err := m.InMe.SetShortURL(ctx, record.ShortURL, record.OriginalURL, record.UserID, record.Deleted); err != nil {
			return fmt.Errorf("error setting short URL in memory: %w", err)
		}
	}

	// Проверяем на ошибки чтения файла
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func (m *InFile) Close() error {
	m.zapLog.Debugln("InFile storage closed (nothing to close currently)")
	return nil
}
