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

const filePermissions = 0o644 // Константа для прав доступа к файлу

// InFile - структура для реализации хранилища в файле.
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

// NewInFile создает новый объект InFile, инициализируя хранилище с поддержкой работы с файлом.
//
// Аргументы:
//   - ctx: Контекст выполнения, который может использоваться для управления временем жизни запроса и отмены.
//   - cfg: Конфигурация, содержащая путь к файлу для хранения данных.
//   - zapLog: Логгер для ведения логов.
//
// Возвращает:
//   - *InFile: Указатель на созданный объект InFile.
//   - error: Ошибка, если она произошла при создании объекта InFile.
//
// Примечание:
// Функция сначала создает хранилище в памяти с помощью inmemory.NewInMemory, а затем читает данные из файла,
// если он существует.
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

// parseShortURLStruct парсит строку JSON в структуру shortURLStruct.
//
// Аргументы:
//   - line: Строка JSON, представляющая запись данных.
//
// Возвращает:
//   - shortURLStruct: Структура, содержащая данные из JSON.
//   - error: Ошибка, если произошла ошибка при парсинге JSON.
//
// Примечание:
// Функция используется для преобразования строки JSON в структуру данных для хранения в памяти.
func parseShortURLStruct(line string) (shortURLStruct, error) {
	var record shortURLStruct
	err := json.Unmarshal([]byte(line), &record)
	if err != nil {
		return shortURLStruct{}, fmt.Errorf("error parsing JSON: %w", err)
	}
	return record, nil
}

// readFileStorage читает данные из файла и загружает их в хранилище в памяти.
//
// Аргументы:
//   - ctx: Контекст выполнения, который может использоваться для управления временем жизни запроса и отмены.
//
// Возвращает:
//   - error: Ошибка, если произошла ошибка при чтении файла или загрузке данных в память.
//
// Примечание:
// Функция открывает файл, считывает данные построчно,
// парсит их и добавляет в хранилище в памяти. Если файл не существует, он создается.
func (f *InFile) readFileStorage(ctx context.Context) error {
	// Открываем файл с флагами для создания, если он не существует
	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, filePermissions)
	if err != nil {
		return fmt.Errorf("error opening or creating file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			f.zapLog.Warnf("error closing file: %w", err)
		}
	}(file)

	// Читаем файл построчно
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Парсим строку JSON в структуру ShortURLStruct
		record, err := parseShortURLStruct(line)
		if err != nil {
			return fmt.Errorf("error parsing line from file: %w", err)
		}

		// Добавляем запись в InMemoryStorage
		if err := f.InMe.SetShortURL(
			ctx,
			record.ShortURL,
			record.OriginalURL,
			record.UserID,
			record.Deleted,
		); err != nil {
			return fmt.Errorf("error setting short URL in memory: %w", err)
		}
	}

	// Проверяем на ошибки чтения файла
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

// Close завершает работу с хранилищем файла. В текущей реализации закрытие файла не требуется,
// поэтому функция просто записывает в лог сообщение о закрытии.
//
// Возвращает:
//   - error: Ошибка, если она произошла при попытке завершить работу хранилища.
func (f *InFile) Close() error {
	f.zapLog.Debugln("InFile storage closed (nothing to close currently)")
	return nil
}
