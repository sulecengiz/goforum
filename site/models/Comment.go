package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Comment modeli
type Comment struct {
	ID         uint `gorm:"primaryKey"`
	PostID     uint
	UserID     *uint
	ParentID   *uint
	Name       string
	Comment    string
	CreatedAt  time.Time
	Replies    []*Comment `gorm:"-"`
	ParentName string     `gorm:"-"`
	LikeCount  int64      `gorm:"-"`
	IsLiked    bool       `gorm:"-"`
}

// Migrasyon
func (c Comment) Migrate() {
	db := GetDB()
	if err := db.AutoMigrate(&c); err != nil {
		fmt.Println("Comment migrate hata:", err)
	}
}

// Ekle
func (c *Comment) AddComment() error { return GetDB().Create(c).Error }

// Hiyerarşik yorumları getir
func (c *Comment) GetCommentTree(postID int) ([]*Comment, error) {
	db := GetDB()
	var comments []*Comment
	if err := db.Where("post_id = ? AND parent_id IS NULL", postID).Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}
	for _, root := range comments {
		if err := c.loadReplies(db, root); err != nil {
			return nil, err
		}
	}
	return comments, nil
}

// Alt yorumları yükle (recursive)
func (c *Comment) loadReplies(db *gorm.DB, parent *Comment) error {
	var replies []*Comment
	if err := db.Where("parent_id = ?", parent.ID).Order("created_at ASC").Find(&replies).Error; err != nil {
		return err
	}
	parent.Replies = replies
	for _, r := range replies {
		if err := c.loadReplies(db, r); err != nil {
			return err
		}
	}
	return nil
}

// Kullanıcının yorumları
func (c Comment) GetUserComments(userID uint) []Comment {
	var list []Comment
	GetDB().Where("user_id = ?", userID).Order("created_at DESC").Find(&list)
	return list
}

// Kullanıcının yorum sayısı
func (c Comment) CountUserComments(userID uint) int64 {
	var cnt int64
	GetDB().Model(&Comment{}).Where("user_id = ?", userID).Count(&cnt)
	return cnt
}

// Yorumu sil
func (c Comment) Delete(commentID uint) error { return GetDB().Delete(&Comment{}, commentID).Error }

func (c Comment) CountByUser(id uint) interface{} {
	var count int64
	GetDB().Model(&Comment{}).Where("user_id = ?", id).Count(&count)
	return count
}
