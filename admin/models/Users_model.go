// models/User.go
package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
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

func (user User) Get(where ...interface{}) User {
	db := GetDB()
	if db == nil {
		return User{}
	}

	var result User
	db.First(&result, where...)
	return result
}

func (user *User) Register() error {
	db := GetDB()
	if db == nil {
		return nil
	}
	return db.Create(user).Error
}

func (user User) Migrate() {
	db := GetDB()
	if db != nil {
		db.AutoMigrate(User{})
	}

}

func (user User) GetAll() interface{} {
	db := GetDB()
	if db == nil {
		return nil
	}

	var users []User
	db.Find(&users)
	return users

}
