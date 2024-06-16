package rwstorage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	PermFile0644 = 0o644
	ErrEmptyFile = "file is empty"
)

type ShortURLStruct struct {
	ID          string `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	Deleted     bool   `json:"deleted"`
}

type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, PermFile0644)
	if err != nil {
		return nil, fmt.Errorf("error open infile %w", err)
	}

	return &Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteData(shData *ShortURLStruct) error {
	data, err := json.Marshal(&shData)
	if err != nil {
		return fmt.Errorf("error marshal data %w", err)
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return fmt.Errorf("error write data %w", err)
	}

	// добавляем перенос строки
	if err := p.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("error write перенос строки %w", err)
	}

	// записываем буфер в файл
	return p.writer.Flush() //nolint: wrapcheck,gocritic // просто пробрасываем ошибку дальше
}

type Consumer struct {
	file *os.File
	// добавляем Reader в Consumer
	Reader *bufio.Scanner
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, PermFile0644)
	if err != nil {
		return nil, fmt.Errorf("error write перенос строки %w", err)
	}

	return &Consumer{
		file: file,
		// создаём новый Reader
		Reader: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) ReadLineInFile() (*ShortURLStruct, error) {
	// читаем данные до символа переноса строки
	data := c.Reader.Text()

	if data == "" {
		return nil, errors.New(ErrEmptyFile)
	}

	// преобразуем данные из JSON-представления в структуру
	shortURL := ShortURLStruct{}
	err := json.Unmarshal([]byte(data), &shortURL)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal shor url: %w", err)
	}

	return &shortURL, nil
}

func (p *Producer) Close() {
	// Закрываем буфер
	if err := p.file.Close(); err != nil {
		log.Fatalf("error producer infile close %v", err)
	}
}

func (c *Consumer) Close() {
	// Закрываем файл
	if err := c.file.Close(); err != nil {
		log.Fatalf("error consumer infile close %v", err)
	}
}

func (p *Producer) WriteBatch(batch []ShortURLStruct) error {
	for _, data := range batch {
		err := p.WriteData(&data)
		if err != nil {
			return fmt.Errorf("ошибка записи данных %w", err)
		}
	}
	return nil
}

func (p *Producer) AppendToFile(sourceFilePath, destFilePath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			log.Printf("error close source file %v", err)
		}
	}(sourceFile)

	destFile, err := os.OpenFile(destFilePath, os.O_APPEND|os.O_WRONLY, PermFile0644)
	if err != nil {
		return fmt.Errorf("error opening destination file: %w", err)
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			log.Printf("error close destination file %v", err)
		}
	}(destFile)

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("error copying data from source to destination: %w", err)
	}

	return nil
}
