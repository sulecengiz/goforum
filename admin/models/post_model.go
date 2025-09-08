package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title, Slug, Description, Content, Picture_url string
	CategoryID                                     int
}

func (p Post) Migrate() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.AutoMigrate(&p)
}

func dbPath() string {
	return "C:/Users/Selly/Desktop/goblog/test.db"
}

func (p Post) Add() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Create(&p)
}
func (p Post) Get(where ...interface{}) Post {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return p
	}
	db.First(&p, where...)
	return p
}
func (p Post) GetAll(where ...interface{}) []Post {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return nil
	}
	var posts []Post
	db.Find(&posts, where...)
	return posts
}
func (p Post) Update(column string, value interface{}) {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Model(&p).Update(column, value)
}
func (p Post) Updates(data Post) {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Model(&p).Updates(data)
}
func (p Post) Delete() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println()
		return
	}
	db.Delete(&p, p.ID)
}
