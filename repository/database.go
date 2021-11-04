package repository

import (
	"errors"
	"github.com/fhmq/hmq/config"
	"github.com/fhmq/hmq/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// IDatabase 数据库接口
type IDatabase interface {
	Orm() *gorm.DB
	Product() IProduct
	Device() IDevice
}

var Database IDatabase

// IProduct 产品数据库接口
type IProduct interface {

	// Add 添加产品
	Add(product *model.Product) error

	// Get 获取产品
	Get(productId string) (*model.Product, error)

	// List 获取产品列表
	List(model.Page) ([]model.Product, model.Page, error)

	// Update 更新产品
	Update(product *model.Product) error

	// Delete 删除指定ID产品
	Delete(productId string) error
}

// IDevice 设备数据库接口
type IDevice interface {

	// Add 添加设备
	Add(device *model.Device) error

	// Get 获取 Device
	Get(productId string, deviceId string) (*model.Device, error)

	// GetSubdevice 获取指定网关对象的子设备
	GetSubdevice(productId string, deviceId string, subProductId string, subDeviceId string) (*model.Device, error)

	// GetConfig 获取设备配置
	GetConfig(productId string, deviceId string) (*string, error)

	// List 获取 Device 列表
	List(page model.Page, device *model.Device) ([]model.Device, model.Page, error)

	// Update 更新 Device
	Update(device *model.Device) error

	// UpdateState 更新设备状态
	UpdateState(productId string, deviceId string, state model.DeviceState) error

	// UpdateConfig 更新设备配置
	UpdateConfig(productId string, deviceId string, config *string) error

	// UpdateGateway 更新网关
	UpdateGateway(productId string, deviceId string, gatewayProductId *string, gatewayDeviceId *string) error

	// UpdateSecret 更新密钥
	UpdateSecret(productId string, deviceId string, deviceSecret string) error

	// Delete 删除指定ID设备
	Delete(productId string, deviceId string) error
}

// Type 数据库类型，目前仅支持两种，mysql和sqlite
type Type string

const (
	MySQLType  Type = "mysql"
	SQLiteType Type = "sqlite"
)

type db struct {
	*gorm.DB
	product IProduct
	device  IDevice
}

func (db *db) Orm() *gorm.DB {
	return db.DB
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
	Database = &db{
		DB:      orm,
		product: &_product{orm},
		device:  &_device{orm},
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
