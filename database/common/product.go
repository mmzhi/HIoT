package common

import (
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/model"
	"gorm.io/gorm"
	"math"
)

type _product struct {
	orm *gorm.DB
}

// Add 添加产品
func (db *_product) Add(product *model.Product) error {
	if tx := db.orm.Create(product); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Get 获取 product
func (db *_product) Get(productId string) (*model.Product, error) {
	var product model.Product
	if tx := db.orm.Where("product_id = ?", productId).First(&product); tx.Error != nil {
		return nil, tx.Error
	}
	return &product, nil
}

// List 获取 product 列表
func (db *_product) List(page model.Page) ([]model.Product, model.Page, error) {
	var products []model.Product
	if tx := db.orm.Model(&model.Product{}).Scopes(database.Paginate(&page)).Find(&products); tx.Error != nil {
		return nil, page, tx.Error
	}
	var total int64
	if tx := db.orm.Model(&model.Product{}).Count(&total); tx.Error != nil {
		return nil, page, tx.Error
	}
	page.Total = int(total)
	page.Pages = int(math.Ceil(float64(page.Total) / float64(page.Size)))
	return products, page, nil
}

// Update 更新 product
func (db *_product) Update(product *model.Product) error {
	if tx := db.orm.Model(product).Select("product_name").Updates(product); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Delete 删除指定ID产品
func (db *_product) Delete(productId string) error {
	if tx := db.orm.Delete(&model.Product{
		ProductId: productId,
	}); tx.Error != nil {
		return tx.Error
	}
	return nil
}
