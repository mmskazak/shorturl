package rwstorage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type RecordURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

const filePerm = 0o600

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, filePerm) // open file
	if err != nil {
		return nil, err //nolint:wrapcheck // Failed to open file
	}

	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteData(shData *RecordURL) error {
	data, err := json.Marshal(&shData)
	if err != nil {
		return err //nolint:wrapcheck // Failed to marshal data
	}

	if _, err := p.writer.Write(data); err != nil {
		return err //nolint:wrapcheck // Failed to write data
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err //nolint:wrapcheck // Failed to write newline
	}

	if err := p.writer.Flush(); err != nil {
		return err //nolint:wrapcheck // Failed to flush buffer
	}

	return nil
}

func (p *Producer) Close() {
	if err := p.file.Close(); err != nil {
		log.Printf("Error closing file: %v", err)
	}
}

type Consumer struct {
	file   *os.File
	reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, filePerm)
	if err != nil {
		return nil, err //nolint:wrapcheck // Failed to open file
	}

	return &Consumer{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadDataFromFile() (*RecordURL, error) {
	data, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err //nolint:wrapcheck // Failed to read data
	}

	event := RecordURL{}
	err = json.Unmarshal(data, &event)
	if err != nil {
		return nil, err //nolint:wrapcheck // Failed to unmarshal data
	}

	return &event, nil
}

func (c *Consumer) Close() {
	if err := c.file.Close(); err != nil {
		log.Printf("Error closing file: %v", err)
	}
}
