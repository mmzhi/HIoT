package repository

import (
	"github.com/ruixiaoedu/hiot/model"
	"gorm.io/gorm"
)

// Paginate 分页方法
func Paginate(page *model.Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.Current <= 0 {
			page.Current = 1
		}
		switch {
		case page.Size > 10000:
			page.Size = 10000
		case page.Size <= 0:
			page.Size = 10
		}
		offset := (page.Current - 1) * page.Size
		return db.Offset(offset).Limit(page.Size)
	}
}
