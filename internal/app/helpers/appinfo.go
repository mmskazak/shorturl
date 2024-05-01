package helpers

import (
	"log"
	"runtime/debug"
	"strings"
)

type AppInfo struct {
	name    string
	version string
}

func GetAppNameAndVersion() AppInfo {
	var app AppInfo

	if info, ok := debug.ReadBuildInfo(); ok {
		app.name = strings.ToUpper(info.Main.Path)

		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				app.version = setting.Value
				break
			}
		}
	}

	if app.name == "" {
		log.Println("name app unknown")
	}
	if app.version == "" {
		log.Println("version app unknown")
	}

	return app
}
