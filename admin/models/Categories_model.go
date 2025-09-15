package models

import "gorm.io/gorm"

// Category admin tarafı modeli
type Category struct {
	gorm.Model
	Title string
	Slug  string
}

func (c Category) Migrate()                          { GetDB().AutoMigrate(&c) }
func (c Category) Add()                              { GetDB().Create(&c) }
func (c Category) Get(where ...interface{}) Category { GetDB().First(&c, where...); return c }
func (c Category) GetAll(where ...interface{}) []Category {
	var list []Category
	GetDB().Find(&list, where...)
	return list
}
func (c Category) Update(column string, value interface{}) { GetDB().Model(&c).Update(column, value) }
func (c Category) Updates(data Category)                   { GetDB().Model(&c).Updates(data) }
func (c Category) Delete()                                 { GetDB().Delete(&c, c.ID) }
