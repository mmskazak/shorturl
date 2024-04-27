package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

var urlMap map[string]string

func main() {
	urlMap = make(map[string]string)

	router := chi.NewRouter()
	router.Get("/", mainPage)
	router.Get("/{id}", handleRedirect)
	router.Post("/", createShortURL)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Сервис сокращения URL"))
	if err != nil {
		return
	}
	return
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Получение значения id из URL-адреса
	id := chi.URLParam(r, "id")

	if len(id) != 8 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL, ok := urlMap[id]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	// Чтение оригинального URL из тела запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read the body", http.StatusBadRequest)
		return
	}
	originalURL := string(body)

	// Генерируем уникальный идентификатор для сокращенной ссылки
	id := generateShortURL(8)
	shortedURL := "http://localhost:8080/" + id
	urlMap[id] = originalURL

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortedURL))
	if err != nil {
		return
	}
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
