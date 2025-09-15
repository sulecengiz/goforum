package models

import (
	"fmt"
	"time"
)

type ContactForm struct {
	// gorm.Model kaldırmadan bırakılabilir ama zaten import edilmiyor burada
	ID                          uint `gorm:"primaryKey"`
	Name, Email, Phone, Message string
	CreatedAt                   string
}

func (c ContactForm) Migrate() { GetDB().AutoMigrate(&c) }

func (c ContactForm) Add() error {
	db := GetDB()
	c.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	if err := db.Create(&c).Error; err != nil {
		fmt.Println("Contact create hata", err)
		return err
	}
	return nil
}

func (c ContactForm) GetAll(where ...interface{}) []ContactForm {
	var contacts []ContactForm
	GetDB().Order("created_at DESC").Find(&contacts, where...)
	return contacts
}

func (c ContactForm) Get(id interface{}) (ContactForm, error) {
	db := GetDB()
	result := db.First(&c, id)
	return c, result.Error
}

func (c ContactForm) Delete() error { return GetDB().Delete(&c).Error }
