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

type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type IGenIDForURL interface {
	Generate() (string, error)
}

type Storage interface {
	Close() error
	GetShortURL(ctx context.Context, id string) (string, error)
	SetShortURL(ctx context.Context, idShortPath string, targetURL string, userID string) error
	SaveBatch(
		ctx context.Context,
		items []Incoming,
		baseHost string,
		userID string,
		generator IGenIDForURL,
	) ([]Output, error)
	GetUserURLs(ctx context.Context, userID string, baseHost string) ([]URL, error)
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
