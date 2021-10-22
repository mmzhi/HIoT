package common

import (
	"github.com/fhmq/hmq/database"
	"gorm.io/gorm"
	"math"
)

type _product struct {
	orm *gorm.DB
}

// Add 添加产品
func (db *_product) Add(product *database.Product) error {
	if tx := db.orm.Create(product); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Get 获取 product
func (db *_product) Get(productId string) (*database.Product, error) {
	var product database.Product
	if tx := db.orm.Where("product_id = ?", productId).First(&product); tx.Error != nil {
		return nil, tx.Error
	}
	return &product, nil
}

// List 获取 product 列表
func (db *_product) List(page database.Page) ([]database.Product, database.Page, error) {
	var products []database.Product
	if tx := db.orm.Model(&database.Product{}).Scopes(database.Paginate(&page)).Find(&products); tx.Error != nil {
		return nil, page, tx.Error
	}
	var total int64
	if tx := db.orm.Model(&database.Product{}).Count(&total); tx.Error != nil {
		return nil, page, tx.Error
	}
	page.Total = int(total)
	page.Pages = int(math.Ceil(float64(page.Total) / float64(page.Size)))
	return products, page, nil
}

// Update 更新 product
func (db *_product) Update(product *database.Product) error {
	if tx := db.orm.Model(product).Select("product_name").Updates(product); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Delete 删除指定ID产品
func (db *_product) Delete(productId string) error {
	return nil
}
