package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("goforum.db"), &gorm.Config{})
	if err != nil {
		panic("Veritabanına bağlanılamadı: " + err.Error())
	}
}

func GetDB() *gorm.DB {
	if DB == nil {
		ConnectDB()
	}
	return DB
}
