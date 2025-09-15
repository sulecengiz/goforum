// site/models/About.go
package models

import (
	"time"
)

// About modeli, sadece bir satır içerecek
type About struct {
	ID        uint   `gorm:"primaryKey;autoIncrement:false"`
	Content   string `gorm:"type:text"`
	UpdatedAt string
}

func (a About) Migrate() {
	db := GetDB()
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
	db := GetDB()
	db.First(&a, 1) // ID 1 olan kaydı getir
	return a, nil
}

func (a About) Update() error {
	db := GetDB()
	a.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return db.Save(&a).Error
}
