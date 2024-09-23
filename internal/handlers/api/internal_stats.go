package api

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/contracts"
	"net/http"
)

func InternalStats(
	ctx context.Context,
	w http.ResponseWriter,
	_ *http.Request,
	store contracts.IInternalStats,
	zapLog *zap.SugaredLogger,
) {
	stats, err := store.InternalStats(ctx)
	if err != nil {
		zapLog.Errorf("Error getting internal stats: %v", err)
	}

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
