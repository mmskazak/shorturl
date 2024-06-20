package api

import (
	"context"
	"encoding/json"
	"log"
	"mmskazak/shorturl/internal/storage"
	"net/http"
)

// DeleteUserURLs - хендлер для асинхронного удаления сокращённых URL по их идентификаторам.
func DeleteUserURLs(ctx context.Context, w http.ResponseWriter, r *http.Request, store storage.Storage) {
	var urlIDs []string
	err := json.NewDecoder(r.Body).Decode(&urlIDs)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received request to delete URLs: %v", urlIDs)

	err = store.DeleteURLs(ctx, urlIDs)
	if err != nil {
		log.Printf("Error deleting URLs: %v", err)
	}

	log.Println("All URLs deletion tasks completed")

	w.WriteHeader(http.StatusAccepted)
}
