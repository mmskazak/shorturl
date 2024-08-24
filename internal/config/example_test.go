package config

import (
	"fmt"
	"time"
)

// ExampleConfig демонстрирует пример базовой(по умолчанию) структуры Config.
func ExampleConfig() {
	config := &Config{
		Address:         ":8080",
		BaseHost:        "http://localhost:8080",
		LogLevel:        "info",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		FileStoragePath: "/tmp/short-url-db.json",
		SecretKey:       "secret",
		DataBaseDSN:     "user:password@/dbname",
	}

	// Форматируем вывод в человекочитаемом виде
	output := fmt.Sprintf(
		`{
	Address: %q,
	BaseHost: %q,
	LogLevel: %q,
	ReadTimeout: %v,
	WriteTimeout: %v,
	FileStoragePath: %q,
	SecretKey: %q,
	DataBaseDSN: %q
}`,
		config.Address,
		config.BaseHost,
		config.LogLevel,
		config.ReadTimeout,
		config.WriteTimeout,
		config.FileStoragePath,
		config.SecretKey,
		config.DataBaseDSN,
	)

	fmt.Println(output)
	// Output:
	// {
	//	Address: ":8080",
	//	BaseHost: "http://localhost:8080",
	//	LogLevel: "info",
	//	ReadTimeout: 10s,
	//	WriteTimeout: 10s,
	//	FileStoragePath: "/tmp/short-url-db.json",
	//	SecretKey: "secret",
	//	DataBaseDSN: "user:password@/dbname"
	// }
}
