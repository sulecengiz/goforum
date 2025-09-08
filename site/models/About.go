// site/models/About.go
package models

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// About modeli, sadece bir satır içerecek
type About struct {
	ID        uint   `gorm:"primaryKey;autoIncrement:false"`
	Content   string `gorm:"type:text"`
	UpdatedAt string
}

func (a About) Migrate() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println("Veritabanı bağlantısı kurulamadı:", err)
		return
	}
	db.AutoMigrate(&a)

	// Eğer tablo boşsa, ilk kaydı oluştur
	var count int64
	db.Model(&a).Count(&count)
	if count == 0 {
		a.ID = 1
		a.Content = "Buraya hakkımda yazısı gelecek."
		a.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		db.Create(&a)
	}
}

func (a About) Get() (About, error) {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		return a, err
	}
	db.First(&a, 1) // ID 1 olan kaydı getir
	return a, nil
}

func (a About) Update() error {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		return err
	}
	a.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	db.Save(&a) // a.ID'ye göre kaydı günceller
	return nil
}
