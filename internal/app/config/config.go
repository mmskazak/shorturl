package config

import (
	"flag"
	"os"
)

type IConfig interface {
	GetAddress() string
	GetBaseHost() string
}

// Config содержит поля вашей конфигурации.
type Config struct {
	address  string
	baseHost string
}

func (c *Config) GetAddress() string {
	return c.address
}

func (c *Config) GetBaseHost() string {
	return c.baseHost
}

func InitConfig() *Config {
	config := &Config{
		address:  ":8080",
		baseHost: "http://localhost:8080",
	}

	// указываем ссылку на переменную, имя флага, значение по умолчанию и описание
	flag.StringVar(&config.address, "a", config.GetAddress(), "Устанавливаем ip адрес нашего сервера.")
	flag.StringVar(&config.baseHost, "b", config.GetBaseHost(), "Устанавливаем базовый URL для сокращенного URL.")

	// делаем разбор командной строки
	flag.Parse()

	// конфигурационные параметры в приоритете из переменных среды
	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		config.address = envServAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.baseHost = envBaseURL
	}

	return config
}
