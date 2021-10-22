package database

import "errors"

// 数据库功能，该插件为必选

// IDatabase 数据库接口
type IDatabase interface {
	Product() IProduct
	Device() IDevice
}

// IProduct 产品数据库接口
type IProduct interface {

	// Add 添加产品
	Add(product *Product) error

	// Get 获取产品
	Get(productId string) (*Product, error)

	// List 获取产品列表
	List(Page) ([]Product, Page, error)

	// Update 更新产品
	Update(product *Product) error

	// Delete 删除指定ID产品
	Delete(productId string) error
}

// IDevice 设备数据库接口
type IDevice interface {

	// Add 添加设备
	Add(device *Device) error

	// Get 获取 Device
	Get(productId string, deviceId string) (*Device, error)

	// GetSubdevice 获取指定网关对象的子设备
	GetSubdevice(productId string, deviceId string, subProductId string, subDeviceId string) (*Device, error)

	// List 获取 Device 列表
	List(page Page) ([]Device, Page, error)

	// Update 更新 Device
	Update(device *Device) error

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

// UnRegister 反注册
func UnRegister(name string) {
	delete(providers, name)
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
