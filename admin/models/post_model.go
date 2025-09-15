package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title, Slug, Description, Content, Picture_url string
	CategoryID                                     int
	AuthorID                                       int        `json:"author_id"`
	Approved                                       bool       `gorm:"default:false" json:"approved"`
	ApprovedAt                                     *time.Time `json:"approved_at"`
}

func (p Post) Migrate() {
	db := GetDB()
	if err := db.AutoMigrate(&p); err != nil {
		fmt.Println("Post migrate hata:", err)
	}
}

func (p Post) Add() {
	GetDB().Create(&p)
}
func (p Post) Get(where ...interface{}) Post {
	GetDB().First(&p, where...)
	return p
}
func (p Post) GetAll(where ...interface{}) []Post {
	var posts []Post
	GetDB().Find(&posts, where...)
	return posts
}
func (p Post) Update(column string, value interface{}) {
	GetDB().Model(&p).Update(column, value)
}
func (p Post) Updates(data Post) {
	GetDB().Model(&p).Updates(data)
}
func (p Post) Delete() {
	GetDB().Delete(&p, p.ID)
}
func (p Post) GetByAuthor(authorID uint) []Post {
	var posts []Post
	GetDB().Where("author_id = ?", authorID).Find(&posts)
	return posts
}
