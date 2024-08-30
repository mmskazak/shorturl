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
	buildDate := time.Now().Format(time.DateTime)

	code := fmt.Sprintf(`package main

var (
	BuildVersion = "%s"
	BuildDate    = "%s"
	BuildCommit  = "%s"
)
`, buildVersion, buildDate, buildCommit)

	file, err := os.Create("version_gen.go")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing file version.go")
		}
	}(file)

	_, err = file.WriteString(code)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
