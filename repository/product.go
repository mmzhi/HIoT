package repository

import (
	"github.com/ruixiaoedu/hiot/model"
	"gorm.io/gorm"
	"math"
)

type _product struct {
	*gorm.DB
}

// Add 添加产品
func (db *_product) Add(product *model.Product) error {
	if tx := db.Create(product); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

// Get 获取 product
func (db *_product) Get(productId string) (*model.Product, error) {
	var product model.Product
	if tx := db.Where("product_id = ?", productId).First(&product); tx.Error != nil {
		return nil, Error(tx.Error)
	}
	return &product, nil
}

// List 获取 product 列表
func (db *_product) List(page model.Page) ([]model.Product, model.Page, error) {
	var products []model.Product
	if tx := db.Model(&model.Product{}).Scopes(Paginate(&page)).Find(&products); tx.Error != nil {
		return nil, page, Error(tx.Error)
	}
	var total int64
	if tx := db.Model(&model.Product{}).Count(&total); tx.Error != nil {
		return nil, page, Error(tx.Error)
	}
	page.Total = int(total)
	page.Pages = int(math.Ceil(float64(page.Total) / float64(page.Size)))
	return products, page, nil
}

// Update 更新 product
func (db *_product) Update(product *model.Product) error {
	if tx := db.Model(product).Select("product_name").Updates(product); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

// Delete 删除指定ID产品
func (db *_product) Delete(productId string) error {
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 先删除设备
		// TODO 删除设备又Device管理
		if tx := tx.Where("product_id = ?", productId).
			Delete(&model.Device{}); tx.Error != nil {
			return tx.Error
		}

		// 再删除产品
		if tx := tx.Delete(&model.Product{
			ProductId: productId,
		}); tx.Error != nil {
			return tx.Error
		}
		return nil
	}); err != nil {
		return Error(err)
	}
	return nil
}
