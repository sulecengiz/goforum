package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Title, Slug string
}

func (p Category) Migrate() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.AutoMigrate(&p)
}

func (p Category) Add() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Create(&p)
}
func (p Category) Get(where ...interface{}) Category {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return p
	}
	db.First(&p, where...)
	return p
}
func (p Category) GetAll(where ...interface{}) []Category {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return nil
	}
	var Categories []Category
	db.Find(&Categories, where...)
	return Categories
}
func (p Category) Update(column string, value interface{}) {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Model(&p).Update(column, value)
}
func (p Category) Updates(data Category) {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Model(&p).Updates(data)
}
func (p Category) Delete() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Delete(&p, p.ID)
}
