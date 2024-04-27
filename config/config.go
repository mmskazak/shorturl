package config

// Config содержит поля вашей конфигурации
type Config struct {
	Address  string
	BaseHost string
}

// NewConfig инициализирует и возвращает новый экземпляр Config с заданными значениями
func CreateConfig() Config {
	return Config{
		Address:  "localhost:8080",
		BaseHost: "localhost:8080",
	}
}
