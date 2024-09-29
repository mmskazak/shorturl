package api

import (
	"context"
	"encoding/json"
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
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	zapLog.Infoln(string(response))
	// Записываем JSON-ответ в ResponseWriter
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		zapLog.Errorf("error writing response: %v", err)
	}
}
