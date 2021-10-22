package common

import (
	"github.com/fhmq/hmq/database"
	"gorm.io/gorm"
)

// CreateDatabase 创建一个数据库对象
func CreateDatabase(orm *gorm.DB) (database.IDatabase, error) {
	err := orm.AutoMigrate(&database.Product{}, &database.Device{})
	if err != nil {
		return nil, err
	}
	return &_db{
		product: &_product{orm},
		device:  &_device{orm},
	}, nil
}

type _db struct {
	product database.IProduct
	device  database.IDevice
}

func (db *_db) Product() database.IProduct {
	return db.product
}

func (db *_db) Device() database.IDevice {
	return db.device
}
