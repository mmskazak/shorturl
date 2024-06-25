package api

import (
	"context"
	"encoding/json"
	"mmskazak/shorturl/internal/storage"
	"net/http"

	"go.uber.org/zap"
)

// DeleteUserURLs - хендлер для асинхронного удаления сокращённых URL по их идентификаторам.
func DeleteUserURLs(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	store storage.Storage,
	zapLog *zap.SugaredLogger,
) {
	var urlIDs []string
	err := json.NewDecoder(r.Body).Decode(&urlIDs)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	zapLog.Infof("Received request to delete URLs: %v", urlIDs)

	err = store.DeleteURLs(ctx, urlIDs)
	if err != nil {
		zapLog.Errorf("Error deleting URLs: %v", err)
	}

	zapLog.Info("All URLs deletion tasks completed")

	w.WriteHeader(http.StatusAccepted)
}
