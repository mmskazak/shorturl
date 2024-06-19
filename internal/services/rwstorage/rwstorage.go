package rwstorage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
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

func (c *Consumer) Close() {
	// Закрываем файл
	if err := c.file.Close(); err != nil {
		log.Fatalf("error consumer infile close %v", err)
	}
}
