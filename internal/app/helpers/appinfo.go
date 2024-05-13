package helpers

import (
	"errors"
	"runtime/debug"
	"strings"
)

type AppInfo struct {
	Name    string
	Version string
}

func GetAppNameAndVersion() (AppInfo, error) {
	var app AppInfo

	if info, ok := debug.ReadBuildInfo(); ok {
		app.Name = strings.ToUpper(info.Main.Path)

		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				app.Version = setting.Value
				break
			}
		}
	} else {
		return app, errors.New("debug.ReadBuildInfo() not have info")
	}

	return app, nil
}
