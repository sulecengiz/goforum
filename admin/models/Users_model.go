package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username, Password string
}

func (p User) Migrate() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.AutoMigrate(&p)
}

func (p User) Add() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Create(&p)
}
func (p User) Get(where ...interface{}) User {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return p
	}
	db.First(&p, where...)
	return p
}
func (p User) GetAll(where ...interface{}) []User {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return nil
	}
	var Users []User
	db.Find(&Users, where...)
	return Users
}
func (p User) Update(column string, value interface{}) {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Model(&p).Update(column, value)
}
func (p User) Updates(data User) {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Model(&p).Updates(data)
}
func (p User) Delete() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Delete(&p, p.ID)
}
