package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var urlMap map[string]string

func main() {
	urlMap = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHortenUrl)
	mux.HandleFunc("/short/", handleRedirect)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

func handleHortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read the body", http.StatusBadRequest)
		return
	}
	originalUrl := string(body)

	// Генерируем уникальный идентификатор для сокращенной ссылки
	shortUrl := generateShortURL(8)
	shortedUrl := "http://localhost:8080/short/" + shortUrl
	urlMap[shortUrl] = originalUrl

	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprint(w, shortedUrl)
	if err != nil {
		return
	}
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortUrl := r.URL.Path[7:]
	originalUrl, ok := urlMap[shortUrl]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusTemporaryRedirect)
}

// generateShortURL генерирует случайный строковый идентификатор заданной длины
func generateShortURL(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	// Создаем новый генератор случайных чисел
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	for i := range b {
		b[i] = charset[rng.Intn(len(charset))] // генерация случайного индекса
	}
	return string(b)
}
