package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"mmskazak/shorturl/internal/contracts"

	"go.uber.org/zap"
)

// InternalStats - внутренняя статистика.
func InternalStats(
	ctx context.Context,
	w http.ResponseWriter,
	_ *http.Request,
	store contracts.IInternalStats,
	zapLog *zap.SugaredLogger,
) {
	zapLog.Info("Getting request by internal stats.")
	stats, err := store.InternalStats(ctx)
	if err != nil {
		zapLog.Errorf("Error getting internal stats: %v", err)
	}
	zapLog.Info("Getting internal stats from store.")

	// Преобразуем данные в JSON
	response, err := json.Marshal(stats)
	if err != nil {
		// Обработка ошибок, связанных с сериализацией JSON
		zapLog.Errorf("error marshalling internal stats: %v", err)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	// Записываем JSON-ответ в ResponseWriter
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		zapLog.Errorf("error writing response: %v", err)
	}
}

// InternalStatsFacade - внутренняя статистика.
func InternalStatsFacade(
	ctx context.Context,
	store contracts.IInternalStats,
	zapLog *zap.SugaredLogger,
) ([]byte, error) {
	zapLog.Info("Getting request by internal stats.")
	stats, err := store.InternalStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting internal stats: %w", err)
	}

	// Преобразуем данные в JSON
	response, err := json.Marshal(stats)
	if err != nil {
		return nil, fmt.Errorf("error marshalling internal stats: %w", err)
	}

	return response, nil
}
