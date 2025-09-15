package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          uint `gorm:"primaryKey;autoIncrement:false"`
	Title, Slug string
}

func (p Category) Migrate()                          { GetDB().AutoMigrate(&p) }
func (p Category) Add()                              { GetDB().Create(&p) }
func (p Category) Get(where ...interface{}) Category { GetDB().First(&p, where...); return p }
func (p Category) GetAll(where ...interface{}) []Category {
	var categories []Category
	GetDB().Find(&categories, where...)
	return categories
}
func (p Category) Update(column string, value interface{}) { GetDB().Model(&p).Update(column, value) }
func (p Category) Updates(data Category)                   { GetDB().Model(&p).Updates(data) }
func (p Category) Delete()                                 { GetDB().Delete(&p, p.ID) }
func (Category) GetBySlug(slug string) Category {
	var c Category
	GetDB().Where("slug = ?", slug).First(&c)
	return c
}
