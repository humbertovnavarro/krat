package db

import (
	"database/sql"
	"embed"
	"fmt"
	"net/url"
	"time"

	"github.com/humbertovnavarro/krat/pkg/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

//go:embed schema.sql
var f embed.FS
var DB *sql.DB
var databaseName = config.NewFilePath(config.AppName + ".db")

var Password = url.QueryEscape(config.GetEnv("DB_PASSWORD"))

func init() {
	logrus.AddHook(NewSqliteHook())
}

func Open() {
	schema, err := f.ReadFile("schema.sql")
	if err != nil {
		logrus.Panic(err)
	}
	connectionString := fmt.Sprintf("file:%s", databaseName)
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		logrus.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		logrus.Fatal(err)
	}
	db.Exec(string(schema))
	DB = db
	Password = ""
	CleanupLogs()
}

func Close() {
	DB.Close()
}

func CleanupLogs() {
	DB.Exec(`DELETE FROM Logs WHERE createdAt < ?`, time.Now().UnixMilli()-DAY.Milliseconds()*7)
}
