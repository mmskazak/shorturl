package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Логика генерации кода
	buildVersion := os.Getenv("BUILD_VERSION")
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	buildCommit := os.Getenv("BUILD_COMMIT")
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	buildDate := os.Getenv("BUILD_DATE")
	if buildDate == "" {
		buildDate = time.Now().Format(time.DateTime)
	}

	code := fmt.Sprintf(`package main

var (
	BuildVersion = "%s"
	BuildDate    = "%s"
	BuildCommit  = "%s"
)
`, buildVersion, buildDate, buildCommit)

	file, err := os.Create("version_gen.go")
	if err != nil {
		log.Printf("Error creating file: %v \n", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file version.go: %v \n", err)
		}
	}(file)

	_, err = file.WriteString(code)
	if err != nil {
		log.Printf("Error writing to file: %v \n", err)
	}
}
