package models

import (
	"gorm.io/gorm"
	"time"
)

type Contact struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"not null"`
	Subject   string         `json:"subject" gorm:"not null"`
	Message   string         `json:"message" gorm:"type:text;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Migrate - Contact tablosunu oluşturma/güncelleme
func (c Contact) Migrate() {
	GetDB().AutoMigrate(&Contact{})
}

// Save - İletişim formunu kaydetme
func (c *Contact) Save() error {
	return GetDB().Create(c).Error
}

// GetAll - Tüm iletişim formlarını getirme
func (c Contact) GetAll() []Contact {
	var contacts []Contact
	GetDB().Order("created_at DESC").Find(&contacts)
	return contacts
}

// Get - ID'ye göre iletişim formu getirme
func (c Contact) Get(id uint) Contact {
	var contact Contact
	GetDB().First(&contact, id)
	return contact
}

// Delete - İletişim formunu silme
func (c Contact) Delete(id uint) error {
	return GetDB().Delete(&Contact{}, id).Error
}
