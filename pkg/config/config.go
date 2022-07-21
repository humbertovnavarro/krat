package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const AppName = "krat"

var AppDataDir string

func init() {
	fetchAppDataDir()
}

func fetchAppDataDir() string {
	d, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatal(err)
	}
	appDir := fmt.Sprintf("%s/.%s", d, AppName)
	AppDataDir = appDir
	return appDir
}

func NewFilePath(directoryStructure ...string) string {
	if len(directoryStructure) == 0 {
		return AppDataDir
	}
	if len(directoryStructure) == 1 {
		return fmt.Sprintf("%s/%s", AppDataDir, directoryStructure[0])
	}
	var path string = directoryStructure[0]
	directoryStructure = directoryStructure[1:]
	for _, name := range directoryStructure {
		path = fmt.Sprintf("%s/%s", path, name)
	}
	return fmt.Sprintf("%s/%s", AppDataDir, path)
}
