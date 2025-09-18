package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type SavedPost struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    uint           `json:"post_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	User User `json:"user" gorm:"foreignKey:UserID"`
	Post Post `json:"post" gorm:"foreignKey:PostID"`
}

// Migrate - SavedPost tablosunu oluşturma/güncelleme
func (sp SavedPost) Migrate() {
	GetDB().AutoMigrate(&SavedPost{})
}

func (sp *SavedPost) Save(userID, postID uint) error {
	db := GetDB()

	// Aynı kullanıcı ve post için zaten kayıtlı mı kontrol et
	var existing SavedPost
	if err := db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existing).Error; err == nil {
		return errors.New("Bu post zaten kaydedilmiş")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	sp.UserID = userID
	sp.PostID = postID
	return db.Create(sp).Error
}

// Unsave - Kaydedilmiş postu kaldırma
func (sp SavedPost) Unsave(userID, postID uint) error {
	return GetDB().Where("user_id = ? AND post_id = ?", userID, postID).Delete(&SavedPost{}).Error
}

// GetByUser - Kullanıcının kaydettiği postları getirme
func (sp SavedPost) GetByUser(userID uint) []SavedPost {
	var savedPosts []SavedPost
	GetDB().Where("user_id = ?", userID).Preload("Post").Find(&savedPosts)
	return savedPosts
}

// IsPostSaved - Post kaydedilmiş mi kontrol et
func (sp SavedPost) IsPostSaved(userID, postID uint) bool {
	var count int64
	GetDB().Model(&SavedPost{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}
