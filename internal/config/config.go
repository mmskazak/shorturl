package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/go-playground/validator/v10"
)

// Config содержит поля конфигурации.
type Config struct {
	Address         string        `validate:"required"`  // Адрес сервера
	BaseHost        string        `validate:"required"`  // Базовый URL
	FileStoragePath string        `validate:"omitempty"` // Путь к файлу хранилища
	DataBaseDSN     string        `validate:"omitempty"` // Строка подключения к базе данных
	SecretKey       string        `validate:"omitempty"` // Секретный ключ для авторизации JWT токена
	LogLevel        LogLevel      `validate:"required"`  // Уровень логирования
	ReadTimeout     time.Duration `validate:"required"`  // Таймаут чтения HTTP-запросов
	WriteTimeout    time.Duration `validate:"required"`  // Таймаут записи HTTP-ответов
}

// validate проверяет правильность заполнения полей конфигурации.
func (c *Config) validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(c)
	if err != nil {
		return fmt.Errorf("error validate %w", err)
	}

	return nil
}

// LogLevel представляет уровень логирования.
type LogLevel string

// Value возвращает уровень логирования в формате zapcore.Level.
func (ll LogLevel) Value() (zapcore.Level, error) {
	switch strings.ToLower(string(ll)) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "dpanic":
		return zapcore.DPanicLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.DebugLevel, errors.New("не найдено соответствие текстовому значению LogLevel, " +
			"уровень логирования задан debug")
	}
}

// InitConfig инициализирует конфигурацию из флагов командной строки и переменных окружения.
// Возвращает указатель на структуру Config и ошибку в случае её возникновения.
func InitConfig() (*Config, error) {
	baseDurationReadTimeout := 10 * time.Second  //nolint:gomnd  // 10 секунд.
	baseDurationWriteTimeout := 10 * time.Second //nolint:gomnd  // 10 секунд.

	config := &Config{
		Address:         ":8080",
		BaseHost:        "http://localhost:8080",
		LogLevel:        "info",
		ReadTimeout:     baseDurationReadTimeout,
		WriteTimeout:    baseDurationWriteTimeout,
		FileStoragePath: "/tmp/short-url-db.json",
		SecretKey:       "secret",
	}

	// Указываем ссылку на переменную, имя флага, значение по умолчанию и описание
	flag.StringVar(&config.Address, "a", config.Address, "IP-адрес сервера")
	flag.StringVar(&config.BaseHost, "b", config.BaseHost, "Базовый URL")
	flag.DurationVar(&config.ReadTimeout, "r", config.ReadTimeout, "ReadTimeout duration")
	flag.DurationVar(&config.WriteTimeout, "w", config.WriteTimeout, "WriteTimeout duration")
	flag.StringVar((*string)(&config.LogLevel), "l", string(config.LogLevel), "log level")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "File storage path")
	flag.StringVar(&config.DataBaseDSN, "d", "", "Database connection string")
	flag.StringVar(&config.SecretKey, "secret", config.SecretKey, "Secret key for authorization JWT token")

	// Разбор командной строки
	flag.Parse()

	// Переопределение значений из переменных окружения, если они заданы
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

	if envLogLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		config.LogLevel = LogLevel(envLogLevel)
	}

	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FileStoragePath = fileStoragePath
	}

	if dataBaseDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.DataBaseDSN = dataBaseDSN
	}

	if secretKey, ok := os.LookupEnv("SECRET_KEY"); ok {
		config.DataBaseDSN = secretKey
	}

	if err := config.validate(); err != nil {
		return &Config{},
			fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	return config, nil
}
