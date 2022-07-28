package db

import (
	"github.com/humbertovnavarro/krat/pkg/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func New(dataDBPath string, logDBPath string) (dataDB *gorm.DB, logDB *gorm.DB) {
	ldb, err := gorm.Open(sqlite.Open(logDBPath), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	db, err := gorm.Open(sqlite.Open(dataDBPath), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	return applySchema(db, ldb)
}
func applySchema(dataDB *gorm.DB, logDB *gorm.DB) (ddb *gorm.DB, ldb *gorm.DB) {
	AutoMigrate(logDB, &models.Log{})
	AutoMigrate(dataDB, &models.OnionService{})
	return dataDB, logDB
}

func AutoMigrate(db *gorm.DB, i interface{}) {
	err := db.AutoMigrate(i)
	if err != nil {
		logrus.Error(err)
	}
}
