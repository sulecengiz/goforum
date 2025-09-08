package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	log.Println("Veritabanı bağlantısı başarılı!")

	// DB.AutoMigrate(&YourModel{})
}
