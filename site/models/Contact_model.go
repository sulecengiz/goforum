package models

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ContactForm struct {
	gorm.Model
	Name, Email, Phone, Message string
	CreatedAt                   string
}

func (c ContactForm) Migrate() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println("Veritabanı bağlantısı kurulamadı:", err)
		return
	}
	db.AutoMigrate(&c)
}

func (c ContactForm) Add() error {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println("Veritabanı bağlantısı kurulamadı:", err)
		return err
	}
	c.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	db.Create(&c)
	return nil
}

func (c ContactForm) GetAll(where ...interface{}) []ContactForm {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println("Veritabanı bağlantısı kurulamadı:", err)
		return nil
	}
	var contacts []ContactForm
	db.Order("created_at DESC").Find(&contacts, where...)
	return contacts
}
func (c ContactForm) Get(id interface{}) (ContactForm, error) {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		return c, err
	}
	result := db.First(&c, id)
	return c, result.Error
}

func (c ContactForm) Delete() error {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		return err
	}
	result := db.Delete(&c)
	return result.Error
}
