package infile

import (
	"context"
	"fmt"
)

// DeleteURLs устанавливает флаг удаления для множества записей.
func (m *InFile) DeleteURLs(ctx context.Context, urlIDs []string) error {
	err := m.InMe.DeleteURLs(ctx, urlIDs)
	if err != nil {
		m.zapLog.Errorf("failed to delete urls: %v", err)
		return fmt.Errorf("delete urls: %w", err)
	}
	m.saveToFile()
	m.zapLog.Infof("deleted %d urls", len(urlIDs))
	return nil
}
