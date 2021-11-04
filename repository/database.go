package repository

import (
	"errors"
	"github.com/ruixiaoedu/hiot/config"
	"github.com/ruixiaoedu/hiot/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// IDatabase 数据库接口
type IDatabase interface {
	Product() IProduct
	Device() IDevice
}

var DB IDatabase

// Type 数据库类型，目前仅支持两种，mysql和sqlite
type Type string

const (
	MySQLType  Type = "mysql"
	SQLiteType Type = "sqlite"
)

type db struct {
	product IProduct
	device  IDevice
}

func (db *db) Product() IProduct {
	return db.product
}

func (db *db) Device() IDevice {
	return db.device
}

// InitDatabase 新建数据库对象
func InitDatabase(cfg config.Database) (err error) {
	var orm *gorm.DB
	switch Type(cfg.Type) {
	case SQLiteType:
	case "":
		if orm, err = initSqlite(cfg); err != nil {
			return err
		}
	case MySQLType:
		if orm, err = initMysql(cfg); err != nil {
			return err
		}
	default:
		return errors.New("unsupported repository type")
	}
	err = orm.AutoMigrate(&model.Product{}, &model.Device{})
	if err != nil {
		return err
	}
	DB = &db{
		product: NewProduct(orm),
		device:  NewDevice(orm),
	}
	return nil
}

// initMysql 初始化sqlite引擎
func initMysql(cfg config.Database) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{
		Logger: NewGormLogger(),
	})
}

// initSqlite 初始化sqlite引擎
func initSqlite(cfg config.Database) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("hiot.db"), &gorm.Config{
		Logger: NewGormLogger(),
	})
}
