package models

import (
	"fmt"
)

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	ParentID  *int      `json:"parent_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Replies   []Comment `json:"replies" gorm:"foreignKey:ParentID"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime"`
}

func (c Comment) Migrate() {
	db := GetDB()
	if err := db.AutoMigrate(&c); err != nil {
		fmt.Println("Comment migrate hata", err)
	}
}

func (c Comment) Add() error {
	return GetDB().Create(&c).Error
}

func (c Comment) GetByPost(postID int) []Comment {
	var comments []Comment
	GetDB().Preload("User").Preload("Replies.User").Where("post_id = ? AND parent_id IS NULL", postID).Find(&comments)
	return comments
}

func (c Comment) Update(data Comment) error {
	return GetDB().Model(&c).Updates(data).Error
}

func (c Comment) Delete() error {
	db := GetDB()
	db.Where("parent_id = ?", c.ID).Delete(&Comment{})
	return db.Delete(&c).Error
}
