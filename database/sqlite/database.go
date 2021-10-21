package mysql

import (
	"github.com/fhmq/hmq/database"
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
	orm, err := gorm.Open(sqlite.Open("hiot.db"), &gorm.Config{
		Logger: database.NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}
	err = orm.AutoMigrate(&database.Product{}, &database.Device{})
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
