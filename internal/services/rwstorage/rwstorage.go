package rwstorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	PermFile0644 = 0o644
)

type ShortURLStruct struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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
	// добавляем reader в Consumer
	reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, PermFile0644)
	if err != nil {
		return nil, fmt.Errorf("error write перенос строки %w", err)
	}

	return &Consumer{
		file: file,
		// создаём новый Reader
		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadDataFromFile() (*ShortURLStruct, error) {
	// читаем данные до символа переноса строки
	data, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err //nolint: wrapcheck,gocritic // нужна чистая ошибка, может быть EOF
	}

	// преобразуем данные из JSON-представления в структуру
	shortURL := ShortURLStruct{}
	err = json.Unmarshal(data, &shortURL)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal shor url %w", err)
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
