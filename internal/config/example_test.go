package config

import (
	"fmt"
	"time"
)

// ExampleConfig демонстрирует пример базовой(по умолчанию) структуры Config.
func ExampleStructConfig() {
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

	fmt.Printf("%+v\n", config)
	// Output:
	// &{Address::8080 BaseHost:http://localhost:8080 FileStoragePath:/tmp/short-url-db.json DataBaseDSN:user:password@/dbname SecretKey:secret LogLevel:info ReadTimeout:10s WriteTimeout:10s}
}
