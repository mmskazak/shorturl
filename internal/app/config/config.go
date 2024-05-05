package config

import (
	"flag"
	"log"
	"os"
	"time"
)

// Config содержит поля вашей конфигурации.
type Config struct {
	Address      string
	BaseHost     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func InitConfig() *Config {
	var baseDurationReadTimeout time.Duration = 10
	var baseDurationWriteTimeout time.Duration = 10

	config := &Config{
		Address:      ":8080",
		BaseHost:     "http://localhost:8080",
		ReadTimeout:  baseDurationReadTimeout,
		WriteTimeout: baseDurationWriteTimeout,
	}

	// указываем ссылку на переменную, имя флага, значение по умолчанию и описание
	flag.StringVar(&config.Address, "a", config.Address, "Устанавливаем ip адрес нашего сервера.")
	flag.StringVar(&config.BaseHost, "b", config.BaseHost, "Устанавливаем базовый URL для сокращенного URL.")
	flag.DurationVar(&config.ReadTimeout, "r", config.ReadTimeout, "ReadTimeout duration.")
	flag.DurationVar(&config.WriteTimeout, "w", config.WriteTimeout, "WriteTimeout duration.")

	// делаем разбор командной строки
	flag.Parse()

	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		config.Address = envServAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.BaseHost = envBaseURL
	}

	if envReadTimeout := os.Getenv("READ_TIMEOUT"); envReadTimeout != "" {
		drt, err := time.ParseDuration(envReadTimeout)
		if err != nil {
			log.Printf("env READ_TIMEOUT не получилось привести к типу \"Duration\": %v", err)
		} else {
			config.ReadTimeout = drt
		}
	}

	if envWriteTimeout := os.Getenv("WRITE_TIMEOUT"); envWriteTimeout != "" {
		dwt, err := time.ParseDuration(envWriteTimeout)
		if err != nil {
			log.Printf("env WRITE_TIMEOUT не получилось привести к типу \"Duration\": %v", err)
		} else {
			config.ReadTimeout = dwt
		}
	}

	return config
}
