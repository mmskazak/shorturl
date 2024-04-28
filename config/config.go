package config

import "flag"

// Config содержит поля вашей конфигурации
type Config struct {
	Address  string
	BaseHost string
}

var config *Config

func InitConfig() *Config {
	config = &Config{
		Address:  ":8080",
		BaseHost: "http://localhost:8080",
	}

	// указываем ссылку на переменную, имя флага, значение по умолчанию и описание
	flag.StringVar(&config.Address, "a", config.Address, "Устанавливаем ip адрес нашего сервера.")
	flag.StringVar(&config.BaseHost, "b", config.BaseHost, "Устанавливаем базовый URL для для сокращенного URL.")

	return config
}

// GetAppConfig CreateConfig NewConfig инициализирует и возвращает новый экземпляр Config с заданными значениями
func GetAppConfig() *Config {
	return config
}
