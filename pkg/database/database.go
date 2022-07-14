package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB
var userConfigDir, _ = os.UserConfigDir()
var UserDir = fmt.Sprintf("%s/tshell", userConfigDir)

func init() {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s", UserDir, "sqlite.db"))
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
