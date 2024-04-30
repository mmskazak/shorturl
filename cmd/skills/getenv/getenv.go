package main

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Files []string `env:"FILES" envSeparator:":"`
	Home  string   `env:"HOME"`
	// required требует, чтобы переменная TASK_DURATION была определена
	TaskDuration time.Duration `env:"TASK_DURATION"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)

	td := os.Getenv("TASK_DURATION")
	fl := os.Getenv("FILES")
	pt := os.Getenv("PATH")

	log.Println("TASK_DURATION: ", td)
	log.Println("FILES: ", fl)
	log.Println("PATH: ", pt)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)
}
