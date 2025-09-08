package models

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Comment modeli
type Comment struct {
	ID         uint `gorm:"primaryKey"`
	PostID     uint
	ParentID   *uint
	Name       string
	Comment    string
	CreatedAt  time.Time
	Replies    []*Comment `gorm:"-"`
	ParentName string     `gorm:"-"`
}

// Comment tablosunu migrate et
func (c Comment) Migrate() {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		fmt.Println("Veritabanı bağlantısı kurulamadı:", err)
		return
	}

	err = db.AutoMigrate(&c)
	if err != nil {
		fmt.Println("Comment tablosu migrate edilirken hata oluştu:", err)
		return
	}

	fmt.Println("Comment tablosu başarıyla migrate edildi.")
}

// Post ID'ye göre tüm yorumları tree şeklinde al
func (c *Comment) GetCommentTree(postID int) ([]*Comment, error) {
	db, _ := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	var allComments []Comment
	db.Where("post_id = ?", postID).Order("created_at asc").Find(&allComments)

	// Map ile hızlı parent-child ilişkisi kur
	commentMap := make(map[uint]*Comment)
	var roots []*Comment

	// Önce tüm comment'leri map'e koy
	for i := range allComments {
		comment := &allComments[i]
		commentMap[comment.ID] = comment
	}

	// Sonra parent-child ilişkilerini kur
	for i := range allComments {
		comment := &allComments[i]
		if comment.ParentID != nil {
			parent := commentMap[*comment.ParentID]
			if parent != nil {
				// ParentName'i set et
				comment.ParentName = parent.Name
				// Child'ı parent'ın replies'ına ekle
				parent.Replies = append(parent.Replies, comment)
			}
		} else {
			// Root comment
			roots = append(roots, comment)
		}
	}

	return roots, nil
}

// Post ID'ye göre yorumları getir
func (c *Comment) GetCommentsByPostID(postID int) ([]Comment, error) {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var comments []Comment
	if err := db.Where("post_id = ?", postID).Order("created_at desc").Find(&comments).Error; err != nil {
		return nil, err
	}

	// Parent isimleri için map oluştur
	commentMap := make(map[uint]string)
	for _, comment := range comments {
		commentMap[comment.ID] = comment.Name
	}

	// Her yoruma parent ismini ekle
	for i := range comments {
		if comments[i].ParentID != nil {
			if parentName, exists := commentMap[*comments[i].ParentID]; exists {
				comments[i].ParentName = parentName
			}
		}
	}

	return comments, nil
}

func (c *Comment) AddComment() error {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		return err
	}
	c.CreatedAt = time.Now()
	result := db.Create(c)
	return result.Error
}
