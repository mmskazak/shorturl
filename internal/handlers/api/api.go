package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/shorturlservice"
	storageInterface "mmskazak/shorturl/internal/storage"
	storageErrors "mmskazak/shorturl/internal/storage/errors"
	"net/http"
)

type JSONRequest struct {
	URL string `json:"url"`
}

type JSONResponse struct {
	ShortURL string `json:"result"`
}

const (
	defaultShortURLLength       = 8
	maxIteration                = 10
	InternalServerErrorMsg      = "Внутренняя ошибка сервера"
	ServiceNotCanCreateShortURL = "Сервису не удалось сформировать короткий URL"
)

func HandleCreateShortURL(w http.ResponseWriter, r *http.Request, storage storageInterface.Storage, baseHost string) {
	// Установка заголовков, чтобы указать, что мы принимаем и отправляем JSON.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")

	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Не удалось прочитать тело запроса %v", err)
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}

	jsonReq := JSONRequest{}

	err = json.Unmarshal(body, &jsonReq)
	if err != nil {
		log.Printf("Ошибка json.Unmarshal: %v", err)
		http.Error(w, ServiceNotCanCreateShortURL, http.StatusInternalServerError)
		return
	}

	generator := genidurl.NewGenIDService()
	shortURLService := shorturlservice.NewShortURLService()
	dto := shorturlservice.DTOShortURL{
		OriginalURL:  jsonReq.URL,
		BaseHost:     baseHost,
		MaxIteration: maxIteration,
		LengthID:     defaultShortURLLength,
	}

	shortURL, err := shortURLService.GenerateShortURL(dto, generator, storage)
	if errors.Is(err, shorturlservice.ErrConflict) {
		shortURLAsJSON, err := buildJSONResponse(shortURL)
		if err != nil {
			log.Printf("Ошибка buildJSONResponse: %v", err)
			http.Error(w, ServiceNotCanCreateShortURL, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte(shortURLAsJSON))
		if err != nil {
			log.Printf("Ошибка write, err := w.Write([]byte(shortURLAsJson)): %v", err)
			http.Error(w, ServiceNotCanCreateShortURL, http.StatusInternalServerError)
			return
		}
		return
	}

	if err != nil {
		log.Printf("Ошибка saveUniqueShortURL: %v", err)
		http.Error(w, ServiceNotCanCreateShortURL, http.StatusInternalServerError)
		return
	}
	jsonResp := JSONResponse{
		ShortURL: shortURL,
	}
	shortURLAsJSON, err := json.Marshal(jsonResp)
	if err != nil {
		log.Printf("Ошибка json.Marshal: %v", err)
		http.Error(w, ServiceNotCanCreateShortURL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(shortURLAsJSON)
	if err != nil {
		log.Printf("Ошибка ResponseWriter: %v", err)
		http.Error(w, InternalServerErrorMsg, http.StatusInternalServerError)
		return
	}
}

func SaveShortenURLsBatch(w http.ResponseWriter, r *http.Request, storage storageInterface.Storage, baseHost string) {
	// Парсинг JSON из тела запроса
	var requestData []storageInterface.Incoming
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	ctx := context.TODO()
	// Сохранение пакета коротких URL
	outputs, err := storage.SaveBatch(ctx, requestData, baseHost)
	if err != nil {
		if errors.Is(err, storageErrors.ErrOriginalURLAlreadyExists) {
			http.Error(w, fmt.Sprintf("conflict saving batch: %v", err), http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("error saving batch: %v", err), http.StatusInternalServerError)
		return
	}

	// Преобразование результата в JSON
	responseData, err := json.Marshal(outputs)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, fmt.Sprintf("error write response body: %v", err), http.StatusInternalServerError)
		return
	}
}

func buildJSONResponse(shortURL string) (string, error) {
	jsonResp := JSONResponse{
		ShortURL: shortURL,
	}
	shortURLAsJSON, err := json.Marshal(jsonResp)
	if err != nil {
		return "", fmt.Errorf("ошибка json marshal %w", err)
	}

	return string(shortURLAsJSON), nil
}
