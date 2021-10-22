package common

import (
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/model"
	"gorm.io/gorm"
)

// 通用的数据库查询，MySQL以及SQLite均在此支持

// CreateDatabase 创建一个数据库对象
func CreateDatabase(orm *gorm.DB) (database.IDatabase, error) {
	err := orm.AutoMigrate(&model.Product{}, &model.Device{})
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
