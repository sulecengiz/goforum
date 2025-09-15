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
	AuthorID                                       uint
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
	db := GetDB()
	db.Create(&p)
}

func (p Post) Get(where ...interface{}) Post {
	db := GetDB()
	db.First(&p, where...)
	return p
}

func (p Post) GetAll(where ...interface{}) []Post {
	db := GetDB()
	var posts []Post
	db.Find(&posts, where...)
	return posts
}

func (p Post) Update(column string, value interface{}) {
	db := GetDB()
	db.Model(&p).Update(column, value)
}

func (p Post) Updates(data Post) {
	db := GetDB()
	db.Model(&p).Updates(data)
}

func (p Post) Delete() {
	db := GetDB()
	db.Delete(&p, p.ID)
}

// Kullanıcının bloglarını getir
func (p Post) GetByAuthor(authorID uint) []Post {
	db := GetDB()
	var posts []Post
	db.Where("author_id = ?", authorID).Find(&posts)
	return posts
}

// Kullanıcının blog sayısını getir
func (p Post) CountByAuthor(authorID uint) int64 {
	db := GetDB()
	var count int64
	db.Model(&Post{}).Where("author_id = ?", authorID).Count(&count)
	return count
}

// Tek bir postu id ve author kontrolü ile getir
func (p Post) GetByIDAndAuthor(id, authorID uint) Post {
	db := GetDB()
	db.Where("id = ? AND author_id = ?", id, authorID).First(&p)
	return p
}

// Post güncelle (sadece ID ile)
func (p Post) EditPost(id uint, updated Post) error {
	db := GetDB()
	return db.Model(&Post{}).Where("id = ?", id).Updates(updated).Error
}

// Post sil (yazar doğrulaması ile)
func (p Post) DeleteByAuthor(id, authorID uint) error {
	db := GetDB()
	return db.Where("id = ? AND author_id = ?", id, authorID).Delete(&Post{}).Error
}

// Genel onaylı ve admin (author_id=0 legacy) postlarını getir
func (p Post) GetPublicPosts() []Post {
	db := GetDB()
	var posts []Post
	db.Where("approved = ? OR author_id = 0", true).Order("created_at DESC").Find(&posts)
	return posts
}

// Admin veya legacy (author_id=0) postları otomatik onayla ve zaman damgası koy
func BackfillAdminPosts() {
	db := GetDB()
	var adminUsers []User
	db.Where("role = ?", 1).Find(&adminUsers)
	var adminIDs []uint
	for _, u := range adminUsers {
		adminIDs = append(adminIDs, u.ID)
	}
	now := time.Now()
	if len(adminIDs) > 0 {
		// Onaylı olmayanları onayla
		db.Model(&Post{}).Where("(author_id IN ? OR author_id = 0) AND approved = ?", adminIDs, false).Updates(map[string]interface{}{"approved": true, "approved_at": &now})
		// Onaylı olup timestamp olmayanları güncelle
		db.Model(&Post{}).Where("(author_id IN ? OR author_id = 0) AND approved = ? AND approved_at IS NULL", adminIDs, true).Update("approved_at", &now)
	} else {
		// Admin yoksa sadece legacy kayıtları işle
		db.Model(&Post{}).Where("author_id = 0 AND approved = ?", false).Updates(map[string]interface{}{"approved": true, "approved_at": &now})
		db.Model(&Post{}).Where("author_id = 0 AND approved = ? AND approved_at IS NULL", true).Update("approved_at", &now)
	}
}
