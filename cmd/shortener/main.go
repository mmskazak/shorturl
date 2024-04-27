package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"mmskazak/shorturl/config"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Добавьте возможность конфигурировать сервис с помощью аргументов командной строки.
// Создайте конфигурацию или переменные для запуска со следующими флагами:
// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
// Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg).
// Совет: создайте отдельный пакет config, где будет храниться структура с вашей конфигурацией и функция, которая будет
// инициализировать поля этой структуры. По мере усложнения конфигурации вы сможете добавлять необходимые поля в вашу
// структуру и инициализировать их.

var urlMap map[string]string

var cfg config.Config

func init() {

	// Создание нового экземпляра конфигурации
	cfg := config.CreateConfig()

	// указываем ссылку на переменную, имя флага, значение по умолчанию и описание
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Устанавливаем ip адрес нашего сервера")
	flag.StringVar(&cfg.BaseHost, "b", cfg.BaseHost, "Устанавливаем ip адрес нашего сервера")
}

func main() {

	// делаем разбор командной строки
	flag.Parse()

	urlMap = make(map[string]string)

	router := chi.NewRouter()
	router.Get("/", mainPage)
	router.Get("/{id}", handleRedirect)
	router.Post("/", createShortURL)

	fmt.Println("Server is running on " + cfg.Address)
	err := http.ListenAndServe(cfg.Address, router)
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
	shortedURL := cfg.BaseHost + "/" + id
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
