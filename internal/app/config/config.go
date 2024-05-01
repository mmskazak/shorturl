package config

import (
	"flag"
	"os"
)

type InitGetConfig interface {
	InitConfig() *Config
	GetAppConfig() *Config
}

// Config содержит поля вашей конфигурации.
type Config struct {
	Address  string
	BaseHost string
}

func InitConfig() *Config {
	config := &Config{
		Address:  ":8080",
		BaseHost: "http://localhost:8080",
	}

	// указываем ссылку на переменную, имя флага, значение по умолчанию и описание
	flag.StringVar(&config.Address, "a", config.Address, "Устанавливаем ip адрес нашего сервера.")
	flag.StringVar(&config.BaseHost, "b", config.BaseHost, "Устанавливаем базовый URL для для сокращенного URL.")

	// делаем разбор командной строки
	flag.Parse()

	// конфигурационные параметры в приоритете из переменных среды
	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		config.Address = envServAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.BaseHost = envBaseURL
	}

	return config
}
