package models

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	ID      int `gorm:"primaryKey;autoIncrement"`
	Level   uint32
	Text    string
	Created int64 `gorm:"autoCreateTime"`
}

type OnionService struct {
	gorm.Model
	ID         string `gorm:"primaryKey"`
	URL        string
	Port       int
	PrivateKey []byte
}

type CryptoKey struct {
	ID         int `gorm:"primaryKey"`
	For        string
	PublicKey  []byte
	PrivateKey []byte
}
