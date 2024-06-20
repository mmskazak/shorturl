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

	// Асинхронное удаление с batch update
	func() {
		batchSize := 10 // Максимальный размер батча
		var batch []string

		for _, id := range urlIDs {
			batch = append(batch, id)

			// Если размер батча достиг максимального, выполняем обновление
			if len(batch) >= batchSize {
				if err := store.DeleteURLs(ctx, batch); err != nil {
					log.Printf("Error deleting batch: %v", err)
				}
				batch = batch[:0] // Очистка батча
			}
		}

		// Обрабатываем оставшиеся записи, если есть
		if len(batch) > 0 {
			if err := store.DeleteURLs(ctx, batch); err != nil {
				log.Printf("Error deleting remaining batch: %v", err)
			}
		}

		log.Println("All URLs deletion tasks completed")
	}()

	w.WriteHeader(http.StatusAccepted)
}
