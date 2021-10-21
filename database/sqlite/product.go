package mysql

import (
	"github.com/fhmq/hmq/database"
	"gorm.io/gorm"
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
func (db *_product) Get(productId string) ( *database.Product, error) {
	var product database.Product
	if tx := db.orm.Where("product_id = ?", productId).First(&product); tx.Error != nil {
		return nil, tx.Error
	}
	return &product, nil
}

// List 获取 product 列表
func (db *_product) List(page int, limit int) ([]database.Product, error) {
	return nil, nil
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