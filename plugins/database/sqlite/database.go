package mysql

import (
	"github.com/fhmq/hmq/plugins/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	err := database.Register(string(database.SQLiteType), &builder{})
	if err != nil {
		return
	}
}

// builder 数据库创建生成器
type builder struct{}

// Build 创建一个数据库对象
func (b *builder) Build(dsn string, extend string) (database.IDatabase, error) {
	db, err := gorm.Open(sqlite.Open("hiot.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&database.Product{}, &database.Device{})
	if err != nil {
		return nil, err
	}
	return &_db{
		product: nil,
		device:  &_device{db},
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
