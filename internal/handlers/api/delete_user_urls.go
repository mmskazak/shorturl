package api

import (
	"context"
	"encoding/json"
	"mmskazak/shorturl/internal/contracts"
	"net/http"

	"go.uber.org/zap"
)

//go:generate mockgen -source=delete_user_urls.go -destination=mocks/mock_delete_user_urls.go -package=mocks

// DeleteUserURLs - хендлер для асинхронного удаления сокращённых URL по их идентификаторам.
// Он принимает JSON-массив с идентификаторами URL, которые нужно удалить.
// Если массив не может быть декодирован или если возникает ошибка при удалении URL,
// он возвращает соответствующий HTTP-статус код.
func DeleteUserURLs(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	store contracts.IDeleteUserURLs,
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
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	zapLog.Info("All URLs deletion tasks completed")

	w.WriteHeader(http.StatusAccepted)
}
