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
	Level      int        `gorm:"-"` // Yorum seviyesi (nested level)
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
		root.Level = 0 // Root comments have level 0
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
	for _, reply := range replies {
		reply.Level = parent.Level + 1 // Set level based on parent
		if err := c.loadReplies(db, reply); err != nil {
			return err
		}
	}
	parent.Replies = replies
	return nil
}

// Yorum sayısını getir
func (c *Comment) GetCommentCount(postID int) (int64, error) {
	var count int64
	if err := GetDB().Model(&Comment{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Yorumu sil
func (c *Comment) DeleteComment(id int) error {
	return GetDB().Delete(&Comment{}, id).Error
}

// Yorum güncelle
func (c *Comment) UpdateComment() error {
	return GetDB().Save(c).Error
}

func (c Comment) GetByUser(id uint) []Comment {
	var comments []Comment
	if err := GetDB().Where("user_id = ?", id).Find(&comments).Error; err != nil {
		return nil
	}
	return comments

}
