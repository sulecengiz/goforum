package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Like modeli - yorumlar için beğeni sistemi
type Like struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"index:idx_user_comment,unique"`
	CommentID uint    `gorm:"index:idx_user_comment,unique"`
	User      User    `gorm:"constraint:OnDelete:CASCADE;"`
	Comment   Comment `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
}

func (l Like) Migrate() {
	db := GetDB()
	if err := db.AutoMigrate(&l); err != nil {
		fmt.Println("Like migrate hata:", err)
		return
	}
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_user_comment ON likes(user_id, comment_id)")
}

func (l Like) AddLike() error {
	db := GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		var cnt int64
		if err := tx.Model(&Like{}).Where("user_id = ? AND comment_id = ?", l.UserID, l.CommentID).Count(&cnt).Error; err != nil {
			return err
		}
		if cnt > 0 {
			// Kaldır
			if err := tx.Where("user_id = ? AND comment_id = ?", l.UserID, l.CommentID).Delete(&Like{}).Error; err != nil {
				return err
			}
			return nil
		}
		// Ekle
		if err := tx.Create(&l).Error; err != nil {
			// Unique / constraint hatasını kullanıcı dostu forma çevir
			if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "constraint") {
				return errors.New("tekrarlı beğeni isteği")
			}
			return err
		}
		return nil
	})
}

// Toggle: Mevcutsa sil (unlike), yoksa ekle (like). liked=true => yeni like eklendi.
func (l Like) Toggle() (liked bool, err error) {
	db := GetDB()
	return liked, db.Transaction(func(tx *gorm.DB) error {
		var existing Like
		errFind := tx.Where("user_id = ? AND comment_id = ?", l.UserID, l.CommentID).First(&existing).Error
		if errFind == nil { // zaten var -> kaldır
			if errDel := tx.Delete(&existing).Error; errDel != nil {
				return errDel
			}
			liked = false
			return nil
		}
		if errFind != nil && !errors.Is(errFind, gorm.ErrRecordNotFound) {
			return errFind
		}
		// yok -> ekle
		if errCreate := tx.Create(&l).Error; errCreate != nil {
			if strings.Contains(errCreate.Error(), "UNIQUE") { // yarışı tolere et
				liked = true
				return nil
			}
			return errCreate
		}
		liked = true
		return nil
	})
}

func (l Like) GetLikeCount(commentID uint) int64 {
	db := GetDB()
	var count int64
	db.Model(&Like{}).Where("comment_id = ?", commentID).Count(&count)
	return count
}

func (l Like) GetLikeCounts(commentIDs []uint) map[uint]int64 {
	res := make(map[uint]int64, len(commentIDs))
	if len(commentIDs) == 0 {
		return res
	}
	db := GetDB()
	type row struct {
		CommentID uint
		Total     int64
	}
	var rows []row
	db.Model(&Like{}).
		Select("comment_id as comment_id, COUNT(*) as total").
		Where("comment_id IN ?", commentIDs).
		Group("comment_id").
		Scan(&rows)
	for _, r := range rows {
		res[r.CommentID] = r.Total
	}
	return res
}

func (l Like) IsLikedByUser(userID, commentID uint) bool {
	if userID == 0 || commentID == 0 {
		return false
	}
	db := GetDB()
	var cnt int64
	db.Model(&Like{}).Where("user_id = ? AND comment_id = ?", userID, commentID).Count(&cnt)
	return cnt > 0
}

func (l Like) GetUserLikedCommentIDs(userID uint, commentIDs []uint) map[uint]bool {
	res := make(map[uint]bool)
	if userID == 0 || len(commentIDs) == 0 {
		return res
	}
	db := GetDB()
	var likes []Like
	db.Where("user_id = ? AND comment_id IN ?", userID, commentIDs).Find(&likes)
	for _, lk := range likes {
		res[lk.CommentID] = true
	}
	return res
}

func (l Like) GetUserLikes(userID uint) []Like {
	db := GetDB()
	var likes []Like
	db.Where("user_id = ?", userID).Find(&likes)
	return likes
}
