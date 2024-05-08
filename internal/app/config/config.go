package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)

// Config содержит поля вашей конфигурации.
type Config struct {
	Address      string        `validate:"required"`
	BaseHost     string        `validate:"required"`
	ReadTimeout  time.Duration `validate:"required"`
	WriteTimeout time.Duration `validate:"required"`
}

func (c Config) validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(c)
	if err != nil {
		return fmt.Errorf("ошибка валидации конфигурации %w", err)
	}

	return nil
}

func InitConfig() (*Config, error) {
	baseDurationReadTimeout := 10 * time.Second  //nolint:gomnd  // 10 секунд.
	baseDurationWriteTimeout := 10 * time.Second //nolint:gomnd  // 10 секунд.

	config := &Config{
		Address:      ":8080",
		BaseHost:     "http://localhost:8080",
		ReadTimeout:  baseDurationReadTimeout,
		WriteTimeout: baseDurationWriteTimeout,
	}

	// указываем ссылку на переменную, имя флага, значение по умолчанию и описание
	flag.StringVar(&config.Address, "a", config.Address, "IP-адерс сервера")
	flag.StringVar(&config.BaseHost, "b", config.BaseHost, "Базовый URL")
	flag.DurationVar(&config.ReadTimeout, "r", config.ReadTimeout, "ReadTimeout duration")
	flag.DurationVar(&config.WriteTimeout, "w", config.WriteTimeout, "WriteTimeout duration")

	// делаем разбор командной строки
	flag.Parse()

	if envServAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.Address = envServAddr
	}

	if envBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.BaseHost = envBaseURL
	}

	if envReadTimeout, ok := os.LookupEnv("READ_TIMEOUT"); ok {
		drt, err := time.ParseDuration(envReadTimeout)
		if err != nil {
			log.Printf("env READ_TIMEOUT не получилось привести к типу \"Duration\": %v", err)
		} else {
			config.ReadTimeout = drt
		}
	}

	if envWriteTimeout, ok := os.LookupEnv("WRITE_TIMEOUT"); ok {
		dwt, err := time.ParseDuration(envWriteTimeout)
		if err != nil {
			log.Printf("env WRITE_TIMEOUT не получилось привести к типу \"Duration\": %v", err)
		} else {
			config.ReadTimeout = dwt
		}
	}

	if err := config.validate(); err != nil {
		return &Config{}, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	return config, nil
}
