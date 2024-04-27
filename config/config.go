package config

// Config содержит поля вашей конфигурации
type Config struct {
	Address  string
	BaseHost string
}

// CreateConfig NewConfig инициализирует и возвращает новый экземпляр Config с заданными значениями
func CreateConfig() Config {
	return Config{
		Address:  ":8080",
		BaseHost: "http://localhost:8080",
	}
}
