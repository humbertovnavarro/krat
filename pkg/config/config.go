package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var UUID string
var AppName string
var Remote string
var AppDataDir string

func init() {
	godotenv.Load()
	AppName = os.Getenv("APP_NAME")
	Remote = os.Getenv("REMOTE")
	fetchAppDataDir()
}

func GetEnv(name string) string {
	return os.Getenv(name)
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
