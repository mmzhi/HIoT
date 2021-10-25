package database

import (
	"errors"
	"github.com/fhmq/hmq/model"
	"gorm.io/gorm"
)

// 数据库功能，该插件为必选

// IDatabase 数据库接口
type IDatabase interface {
	Product() IProduct
	Device() IDevice
}

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

	// Delete 删除指定ID设备
	Delete(productId string, deviceId string) error
}

// Type 数据库类型，目前仅支持两种，mysql和sqlite
type Type string

const (
	MySQLType  Type = "mysql"
	SQLiteType Type = "sqlite"
)

// IBuilder 构建器
type IBuilder interface {
	// Build DSN信息，可参照https://gorm.io/zh_CN/docs/connecting_to_the_database.html
	Build(dsn string, extend string) (IDatabase, error)
}

var providers = make(map[string]IBuilder)

// Register 数据库方法
func Register(name string, i IBuilder) error {
	if name == "" || i == nil {
		return errors.New("invalid args")
	}

	if _, dup := providers[name]; dup {
		return errors.New("already exists")
	}

	providers[name] = i

	return nil
}

var _database IDatabase

// Database 获取数据库
func Database() (IDatabase, error) {
	if _database == nil {
		return nil, errors.New("database is not initialization")
	}
	return _database, nil
}

// InitDatabase 新建数据库对象
func InitDatabase(name string, dsn string, extend string) (err error) {
	if name == "" {
		name = string(SQLiteType)
	}
	builder, ok := providers[name]
	if !ok {
		return errors.New("not exists")
	}
	_database, err = builder.Build(dsn, extend)
	return err
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
