package storage

import (
	"context"
	"fmt"
	"net/url"
)

type Incoming struct {
	CorrelationID string `json:"correlation_id"` // строковый идентификатор
	OriginalURL   string `json:"original_url"`   // оригинальный URL
}

type Output struct {
	CorrelationID string `json:"correlation_id"` // строковый идентификатор
	ShortURL      string `json:"short_url"`      // короткий URL
}

type Storage interface {
	Close() error
	GetShortURL(ctx context.Context, id string) (string, error)
	SetShortURL(ctx context.Context, id string, targetURL string) error
	SaveBatch(ctx context.Context, items []Incoming, baseHost string) ([]Output, error)
}

func GetFullShortURL(baseHost, correlationID string) (string, error) {
	u, err := url.Parse(baseHost)
	if err != nil {
		return "", fmt.Errorf("error parsing baseHost: %w", err)
	}
	// ResolveReference correctly concatenates the base URL and the path
	u = u.ResolveReference(&url.URL{Path: correlationID})
	return u.String(), nil
}
