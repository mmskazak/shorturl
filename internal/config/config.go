package config

import (
	"encoding/json"
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
	Address         string        `json:"address" validate:"required"`            // Адрес сервера
	BaseHost        string        `json:"base_host" validate:"required"`          // Базовый URL
	FileStoragePath string        `json:"file_storage_path" validate:"omitempty"` // Путь к файлу хранилища
	DataBaseDSN     string        `json:"database_dsn" validate:"omitempty"`      // Строка подключения к базе данных
	SecretKey       string        `json:"secret_key" validate:"omitempty"`        // Секретный ключ JWT токена
	ConfigPath      string        `json:"config_path" validate:"omitempty"`       // Путь к конфигурационному файлу
	LogLevel        LogLevel      `json:"log_level" validate:"required"`          // Уровень логирования
	ReadTimeout     time.Duration `json:"read_timeout" validate:"required"`       // Таймаут чтения HTTP-запросов
	WriteTimeout    time.Duration `json:"write_timeout" validate:"required"`      // Таймаут записи HTTP-ответов
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
		ConfigPath:      "",
	}

	// Указываем ссылку на переменную, имя флага, значение по умолчанию и описание
	flag.StringVar(&config.ConfigPath, "c", config.ConfigPath, "Path to configuration file")
	flag.StringVar(&config.ConfigPath, "config", config.ConfigPath, "Path to configuration file")
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
	if envConfigPath, ok := os.LookupEnv("CONFIG"); ok {
		config.ConfigPath = envConfigPath
	}

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
			config.WriteTimeout = dwt
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
		config.SecretKey = secretKey
	}

	if err := config.validate(); err != nil {
		return &Config{}, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	// Конфигурационный файл имеет наименьший приоритет

	if config.ConfigPath != "" {
		// Читаем конфиг
		data, err := os.ReadFile(config.ConfigPath)
		if err != nil {
			log.Fatalf("Ошибка при чтении файла: %v", err)
		}

		// Инициализируем структуру конфигурации
		var configFromFile map[string]interface{}

		// Декодируем JSON данные в структуру
		err = json.Unmarshal(data, &configFromFile)
		if err != nil {
			log.Fatalf("Ошибка при парсинге JSON: %v", err)
		}

		// Переносим значения из файла в основную конфигурацию, если они не установлены
		assignConfigDefaults(config, configFromFile)
	}

	return config, nil
}

// assignConfigDefaults - функция для переноса значений из второго конфига в основной,
// если в основном они равны значениям по умолчанию.
func assignConfigDefaults(config *Config, configFromFile map[string]interface{}) {
	if config.Address == ":8080" {
		if addr, ok := configFromFile["address"].(string); ok && addr != "" {
			config.Address = addr
		}
	}
	if config.BaseHost == "http://localhost:8080" {
		if baseHost, ok := configFromFile["base_host"].(string); ok && baseHost != "" {
			config.BaseHost = baseHost
		}
	}
	if config.LogLevel == "info" {
		if logLevel, ok := configFromFile["log_level"].(string); ok && logLevel != "" {
			config.LogLevel = LogLevel(logLevel)
		}
	}
	if config.ReadTimeout == 10*time.Second {
		// Пример обработки `read_timeout`, который может быть представлен как строка или число
		if readTimeoutStr, ok := configFromFile["read_timeout"].(string); ok && readTimeoutStr != "" {
			// Обработка строки
			if duration, err := time.ParseDuration(readTimeoutStr); err == nil {
				config.ReadTimeout = duration
			}
		} else if readTimeoutNum, ok := configFromFile["read_timeout"].(float64); ok {
			// Обработка числа (в секундах)
			config.ReadTimeout = time.Duration(readTimeoutNum) * time.Second
		}
	}
	// Если значение write_timeout в конфигурации по умолчанию
	if config.WriteTimeout == 10*time.Second {
		// Обработка, если значение read_timeout представлено как строка
		if writeTimeoutStr, ok := configFromFile["write_timeout"].(string); ok && writeTimeoutStr != "" {
			if duration, err := time.ParseDuration(writeTimeoutStr); err == nil {
				config.WriteTimeout = duration
			}
		} else if writeTimeoutNum, ok := configFromFile["write_timeout"].(float64); ok {
			// Обработка, если значение write_timeout представлено как число
			config.WriteTimeout = time.Duration(writeTimeoutNum) * time.Second
		}
	}
	if config.FileStoragePath == "/tmp/short-url-db.json" {
		if fileStoragePath, ok := configFromFile["file_storage_path"].(string); ok && fileStoragePath != "" {
			config.FileStoragePath = fileStoragePath
		}
	}
	if config.DataBaseDSN == "" {
		if dataBaseDSN, ok := configFromFile["database_dsn"].(string); ok && dataBaseDSN != "" {
			config.DataBaseDSN = dataBaseDSN
		}
	}
	if config.SecretKey == "secret" {
		if secretKey, ok := configFromFile["secret_key"].(string); ok && secretKey != "" {
			config.SecretKey = secretKey
		}
	}
}
