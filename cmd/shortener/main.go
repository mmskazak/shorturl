package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"mmskazak/shorturl/config"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

var urlMap map[string]string
var cfg config.Config

func init() {
	// Создание нового экземпляра конфигурации
	cfg = config.CreateConfig()
}

func main() {
	// указываем ссылку на переменную, имя флага, значение по умолчанию и описание
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Устанавливаем ip адрес нашего сервера.")
	flag.StringVar(&cfg.BaseHost, "b", cfg.BaseHost, "Устанавливаем базовый URL для для сокращенного URL.")

	//конфигурационные параметры в приоритете из переменных среды
	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		cfg.Address = envServAddr
	}
	if envBaseUrl := os.Getenv("BASE_URL"); envBaseUrl != "" {
		cfg.BaseHost = envBaseUrl
	}

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
