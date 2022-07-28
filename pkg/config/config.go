package config

import (
	"fmt"
	"os"

	"github.com/humbertovnavarro/krat/pkg/db"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var UUID string
var AppName string
var Remote string
var AppDataDir string
var DB *gorm.DB
var LogDB *gorm.DB

func init() {
	godotenv.Load()
	AppName = os.Getenv("APP_NAME")
	Remote = os.Getenv("REMOTE")
	d, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatal(err)
	}
	AppDataDir = d
	db, ldb := db.New(FilePath("data.db"), FilePath("log.db"))
	logrus.Infof("opened database files @ %s", FilePath())
	DB = db
	LogDB = ldb
}

func GetEnv(name string) string {
	return os.Getenv(name)
}

// returns the file directory structure specified by variable args in app storage directory
func FilePath(directoryStructure ...string) string {
	if len(directoryStructure) == 0 {
		return fmt.Sprintf("%s/.%s", AppDataDir, AppName)
	}
	if len(directoryStructure) == 1 {
		return fmt.Sprintf("%s/.%s/%s", AppDataDir, AppName, directoryStructure[0])
	}
	var path string = directoryStructure[0]
	directoryStructure = directoryStructure[1:]
	for _, name := range directoryStructure {
		path = fmt.Sprintf("%s/%s", path, name)
	}
	return fmt.Sprintf("%s/.%s/%s", AppDataDir, AppName, path)
}
