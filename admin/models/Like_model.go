package models

import (
	"fmt"
)

type Like struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	PostID    int  `json:"post_id"`
	CommentID *int `json:"comment_id,omitempty"`
}

func (l Like) Migrate() {
	if err := GetDB().AutoMigrate(&l); err != nil {
		fmt.Println("Like migrate hata", err)
	}
}
func (l Like) Add() error {
	return GetDB().Create(&l).Error
}
func (l Like) Delete() error {
	return GetDB().Delete(&l).Error
}
func (l Like) Exists() bool {
	var count int64
	GetDB().Model(&Like{}).Where("user_id = ? AND post_id = ?", l.UserID, l.PostID).Count(&count)
	return count > 0
}
