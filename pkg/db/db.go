package db

import (
	"database/sql"
	"embed"

	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/sirupsen/logrus"
)

//go:embed schema.sql
var f embed.FS
var DB *sql.DB

func init() {
	schema, err := f.ReadFile("schema.sql")
	if err != nil {
		logrus.Panic(err)
	}
	sqlFile := config.NewFilePath("db.sqlite")
	db, err := sql.Open("sqlite3", sqlFile)
	if err != nil {
		logrus.Fatal(err)
	}
	db.Exec(string(schema))
	DB = db
}

func Close() {
	DB.Close()
}
