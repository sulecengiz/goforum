package models

import (
	"time"

	"gorm.io/gorm"
)

type SavedPost struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	PostID    uint           `json:"post_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

func (SavedPost) TableName() string {
	return "saved_posts"
}

func (s SavedPost) Migrate() {
	DB.AutoMigrate(&SavedPost{})
}

// Create - Yeni kaydetme işlemi
func (s SavedPost) Create() SavedPost {
	DB.Create(&s)
	return s
}

// Get - Tek kayıt getir
func (s SavedPost) Get(where ...interface{}) SavedPost {
	DB.Where(where[0], where[1:]...).First(&s)
	return s
}

// GetAll - Tüm kayıtları getir
func (s SavedPost) GetAll(where ...interface{}) []SavedPost {
	var savedPosts []SavedPost
	query := DB
	if len(where) > 0 {
		query = query.Where(where[0], where[1:]...)
	}
	query.Find(&savedPosts)
	return savedPosts
}

// Delete - Kayıt sil
func (s SavedPost) Delete() {
	DB.Delete(&s)
}

// IsPostSavedByUser - Kullanıcının postu kaydetmiş mi kontrol et
func (s SavedPost) IsPostSavedByUser(userID, postID uint) bool {
	var count int64
	DB.Model(&SavedPost{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

// GetSavedPostsByUser - Kullanıcının kaydettiği postları getir
func (s SavedPost) GetSavedPostsByUser(userID uint) []SavedPost {
	var savedPosts []SavedPost
	DB.Preload("Post").Preload("Post.Category").Where("user_id = ?", userID).Order("created_at DESC").Find(&savedPosts)
	return savedPosts
}

// GetSavedPostIDs - Kullanıcının kaydettiği post ID'lerini getir
func (s SavedPost) GetSavedPostIDs(userID uint) []uint {
	var postIDs []uint
	DB.Model(&SavedPost{}).Where("user_id = ?", userID).Pluck("post_id", &postIDs)
	return postIDs
}

// ToggleSave - Kaydetme durumunu değiştir (kaydet/kaldır)
func (s SavedPost) ToggleSave(userID, postID uint) (bool, error) {
	existing := SavedPost{}.Get("user_id = ? AND post_id = ?", userID, postID)

	if existing.ID != 0 {
		// Kayıt varsa sil
		existing.Delete()
		return false, nil
	} else {
		// Kayıt yoksa oluştur
		newSave := SavedPost{
			UserID: userID,
			PostID: postID,
		}
		newSave.Create()
		return true, nil
	}
}

// GetByUser - Kullanıcının kaydettiği postları getir (eski metodla uyumlu)
func (s SavedPost) GetByUser(userID uint) []SavedPost {
	return s.GetSavedPostsByUser(userID)
}
