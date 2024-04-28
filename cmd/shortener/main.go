package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"mmskazak/shorturl/config"
	"mmskazak/shorturl/internal/handlers"
	"mmskazak/shorturl/internal/repository"
	"net/http"
	"os"
)

var urlMap map[string]string
var cfg *config.Config

func init() {
	// Создание нового экземпляра конфигурации
	cfg = config.InitConfig()
	repository.InitUrlMap()
}

func main() {

	// делаем разбор командной строки
	flag.Parse()

	//конфигурационные параметры в приоритете из переменных среды
	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		cfg.Address = envServAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseHost = envBaseURL
	}

	urlMap = make(map[string]string)

	router := chi.NewRouter()
	router.Get("/", mainPage)
	router.Get("/{id}", handlers.HandleRedirect)
	router.Post("/", handlers.CreateShortURL)

	fmt.Println("Server is running on " + cfg.Address)
	err := http.ListenAndServe(cfg.Address, router)
	if err != nil {
		return
	}
}

func mainPage(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Сервис сокращения URL"))
	if err != nil {
		return
	}
}
