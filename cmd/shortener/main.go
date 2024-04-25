package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urlMap map[string]string

func main() {
	urlMap = make(map[string]string)

	router := http.NewServeMux()
	router.HandleFunc("/", handleRequests)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// Разделите путь на компоненты
	parts := strings.Split(path, "/")
	// Проверяем, есть ли после базового '/' еще один компонент пути
	if len(parts) >= 2 && parts[1] != "" {
		handleRedirect(w, r)
	} else {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Сервис сокращения URL"))
			return
		}
		handleHortenURL(w, r)
	}
}

func handleHortenURL(w http.ResponseWriter, r *http.Request) {
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
	originalURL := string(body)

	// Генерируем уникальный идентификатор для сокращенной ссылки
	shortURL := generateShortURL(8)
	shortedURL := "http://localhost:8080/" + shortURL
	urlMap[shortURL] = originalURL

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortedURL))
	if err != nil {
		return
	}
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortURL := r.URL.Path[1:]
	if len(shortURL) != 8 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL, ok := urlMap[shortURL]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
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
