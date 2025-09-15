package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
}

func (user *User) Login() *User {
	db := GetDB()
	if db == nil {
		return nil
	}

	var result User
	if err := db.Where("username = ? AND password = ?", user.Username, user.Password).First(&result).Error; err != nil {
		return nil
	}
	return &result
}

func (user *User) Register() error {
	db := GetDB()
	if db == nil {
		return errors.New("veritabanı bağlantısı kurulamadı")
	}

	// Kullanıcı adı kontrolü
	var existingUser User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err != gorm.ErrRecordNotFound {
		return errors.New("bu kullanıcı adı zaten kullanımda")
	}

	// E-posta kontrolü
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != gorm.ErrRecordNotFound {
		return errors.New("bu e-posta adresi zaten kullanımda")
	}

	// Yeni kullanıcıyı kaydet
	if err := db.Create(user).Error; err != nil {
		return errors.New("kullanıcı kaydedilirken bir hata oluştu: " + err.Error())
	}

	return nil
}

func (user User) Get(id uint) User {
	db := GetDB()
	if db == nil {
		return User{}
	}
	var result User
	db.First(&result, id)
	return result
}
