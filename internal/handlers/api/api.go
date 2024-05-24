package api

import (
	"encoding/json"
	"io"
	"mmskazak/shorturl/internal/logger"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/shorturlservice"
	"net/http"
)

type IStorage interface {
	GetShortURL(id string) (string, error)
	SetShortURL(id string, targetURL string) error
}

const (
	defaultShortURLLength       = 8
	maxIteration                = 10
	InternalServerErrorMsg      = "Внутренняя ошибка сервера"
	ServiceNotCanCreateShortURL = "Сервису не удалось сформировать короткий URL"
)

func HandleCreateShortURL(w http.ResponseWriter, r *http.Request, storage IStorage, baseHost string) {
	// Установка заголовков, чтобы указать, что мы принимаем и отправляем JSON.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")

	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Logf.Errorf("Не удалось прочитать тело запроса %v", err)
		http.Error(w, "Что-то пошло не так!", http.StatusBadRequest)
		return
	}

	type JSONRequest struct {
		URL string `json:"url"`
	}
	jsonReq := JSONRequest{}

	err = json.Unmarshal(body, &jsonReq)
	if err != nil {
		logger.Logf.Errorf("Ошибка json.Unmarshal: %v", err)
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
	if err != nil {
		logger.Logf.Errorf("Ошибка saveUniqueShortURL: %v", err)
		http.Error(w, ServiceNotCanCreateShortURL, http.StatusInternalServerError)
		return
	}

	type JSONResponse struct {
		ShortURL string `json:"result"`
	}

	jsonResp := JSONResponse{
		ShortURL: shortURL,
	}
	shortURLAsJSON, err := json.Marshal(jsonResp)
	if err != nil {
		logger.Logf.Errorf("Ошибка json.Marshal: %v", err)
		http.Error(w, ServiceNotCanCreateShortURL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(shortURLAsJSON)
	if err != nil {
		logger.Logf.Errorf("Ошибка ResponseWriter: %v", err)
		http.Error(w, InternalServerErrorMsg, http.StatusInternalServerError)
		return
	}
}
